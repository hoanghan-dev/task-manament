package services

import (
	"context"
	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/dto/response"
	"dev/task-management/internal/modules/task/events"
	"dev/task-management/internal/modules/task/mapper"
	"dev/task-management/internal/modules/task/repositories"
	workspaceService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/pkg/apperror"
	"dev/task-management/pkg/cache"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
	GetAllTask(ctx context.Context, ownerId uuid.UUID) ([]*response.TaskResponse, error)
	GetTask(context context.Context, id uuid.UUID, ownerId uuid.UUID) (*response.TaskResponse, error)
	CreateTask(ctx context.Context, t *request.CreateTaskRequest, ownerId uuid.UUID) (*response.TaskResponse, error)
	UpdateTask(ctx context.Context, id uuid.UUID, t *request.UpdateTaskRequest, ownerId uuid.UUID) error
	DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error
	AssignTask(ctx context.Context, assignTaskReq *request.AssignTaskRequest, ownerId uuid.UUID) error
	UpdateTaskStatus(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, taskReq *request.UpdateTaskStatusRequest) error
}

type taskService struct {
	taskRepository   repositories.TaskRepository
	workspaceService workspaceService.WorkspaceService
	redisClient      *cache.RedisCacheService
}

func NewTaskService(repo repositories.TaskRepository,
	workspaceService workspaceService.WorkspaceService,
	client *cache.RedisCacheService) TaskService {
	return &taskService{
		taskRepository:   repo,
		workspaceService: workspaceService,
		redisClient:      client,
	}
}

func (s *taskService) GetAllTask(ctx context.Context, ownerId uuid.UUID) ([]*response.TaskResponse, error) {

	redisKey := "tasks:owner_id:" + ownerId.String()

	tasksCache := make([]*response.TaskResponse, 0)

	err := s.redisClient.Get(redisKey, &tasksCache)

	if err == nil {
		fmt.Println("[INFO] Jump into cache")
		return tasksCache, nil
	}
	fmt.Println("[INFO] Get data into database")
	tasks, err := s.taskRepository.FindAll(ctx, ownerId)

	if err != nil {
		return nil, err // AppError from repository
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	taskRes := make([]*response.TaskResponse, 0)

	for _, task := range tasks {
		t := mapper.EntityToTaskResponse(&task)
		taskRes = append(taskRes, t)
	}

	errRedis := s.redisClient.Set(redisKey, taskRes, 10)

	if errRedis != nil {
		fmt.Printf("[ERROR] Set task list into cache failed with error: %v\n", errRedis)
	}
	return taskRes, nil
}

func (s *taskService) GetTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) (*response.TaskResponse, error) {

	redisKey := "tasks:owner_id:" + ownerId.String() + ":task_id:" + id.String()

	taskCache := &response.TaskResponse{}

	err := s.redisClient.Get(redisKey, &taskCache)

	if err == nil {
		return taskCache, nil
	}

	task, err := s.taskRepository.FindById(ctx, id, ownerId)
	if err != nil {
		// Repository returns AppError (NotFound or Internal) → propagate
		return nil, err
	}

	taskRes := mapper.EntityToTaskResponse(task)

	// Cache set failure is non-critical → just log
	if cacheErr := s.redisClient.Set(redisKey, taskRes, 10); cacheErr != nil {
		fmt.Printf("[WARN] Failed to cache task: %v\n", cacheErr)
	}

	return taskRes, nil
}

func (s *taskService) CreateTask(ctx context.Context, t *request.CreateTaskRequest, ownerId uuid.UUID) (*response.TaskResponse, error) {

	ownered, err := s.workspaceService.IsWorkspaceOwnedBy(ctx, t.Workspace, ownerId)

	if err != nil {
		return nil, err // AppError from workspace service
	}

	if !ownered {
		return nil, apperror.NewForbidden("you don't have permission to create task in this workspace")
	}

	id := uuid.New()
	createAt := time.Now()
	task := mapper.TaskRequestToEntity(id, t, createAt)

	if !task.StatusIsValid() {
		return nil, apperror.NewValidation("task status invalid, must be one of: TODO, IN_PROGRESS, DONE, BLOCKED")
	}

	createErr := s.taskRepository.CreateTask(ctx, task)

	if createErr != nil {
		return nil, createErr // AppError from repository
	}

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] Clear redis task cache failed with error: %v\n", errRedis)
	}

	return mapper.EntityToTaskResponse(task), nil
}

func (s *taskService) UpdateTask(ctx context.Context, id uuid.UUID, t *request.UpdateTaskRequest, ownerId uuid.UUID) error {
	task := mapper.TaskUpdateToEntity(t)
	err := s.taskRepository.UpdateTask(ctx, id, task, ownerId)
	if err != nil {
		return err // AppError from repository (NotFound or Internal)
	}

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] Clear redis task cache failed with error: %v\n", errRedis)
	}
	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	err := s.taskRepository.DeleteTask(ctx, id, ownerId)
	if err != nil {
		return err // AppError from repository (NotFound or Internal)
	}

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] Clear redis task cache failed with error: %v\n", errRedis)
	}
	return nil
}

func (s *taskService) AssignTask(ctx context.Context, assignTaskReq *request.AssignTaskRequest, ownerId uuid.UUID) error {
	if !s.taskRepository.TaskExistsInWorkspace(ctx, assignTaskReq.TaskId, assignTaskReq.WorkspaceId) {
		return apperror.NewNotFound("task in workspace")
	}

	exists := s.taskRepository.TaskIsExists(ctx, assignTaskReq.TaskId)

	if !exists {
		return apperror.NewNotFound("task")
	}

	ownered, err := s.workspaceService.IsWorkspaceOwnedBy(ctx, assignTaskReq.WorkspaceId, ownerId)

	if err != nil {
		return err // AppError from workspace service
	}

	if !ownered {
		return apperror.NewForbidden("you don't have permission to assign tasks in this workspace")
	}

	errAssign := s.taskRepository.AssignTask(ctx, assignTaskReq.TaskId, assignTaskReq.AssigneeId)

	if errAssign != nil {
		return errAssign // AppError from repository (Conflict or Internal)
	}

	go s.notifyAssignTask(ownerId, assignTaskReq.AssigneeId, assignTaskReq.TaskId)

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] Clear redis task cache failed with error: %v\n", errRedis)
	}
	return nil
}
func (s *taskService) notifyAssignTask(senderId uuid.UUID, receiverId uuid.UUID, taskId uuid.UUID) {
	queueKey := "queue:notifications:assign"

	event := events.NewAssignTaskEvent(senderId, receiverId, taskId)

	err := s.redisClient.Push(context.Background(), queueKey, event)

	if err != nil {
		fmt.Printf("[ERROR] publish task assigned event failed: %v\n", err)
	}

}

func (s *taskService) UpdateTaskStatus(ctx context.Context, taskId uuid.UUID, ownerId uuid.UUID, taskReq *request.UpdateTaskStatusRequest) error {
	if !s.taskRepository.TaskExistsInWorkspace(ctx, taskId, taskReq.WorkspaceId) {
		return apperror.NewNotFound("task in workspace")
	}

	task, err := s.taskRepository.FindById(ctx, taskId, ownerId)

	if err != nil {
		return err // AppError from repository (NotFound or Internal)
	}

	ownered, err := s.workspaceService.IsWorkspaceOwnedBy(ctx, taskReq.WorkspaceId, ownerId)

	if err != nil {
		return err // AppError from workspace service
	}

	if !ownered {
		return apperror.NewForbidden("you don't have permission to update task status in this workspace")
	}

	if !s.changeStatusIsValid(task.Status, taskReq.Status) {
		return apperror.NewBadRequest(fmt.Sprintf("cannot change task status from %s to %s", task.Status, taskReq.Status))
	}

	err = s.taskRepository.UpdateTaskStatus(ctx, taskId, taskReq.Status)

	if err != nil {
		return err // AppError from repository
	}

	s.notifyUpdateTaskStatus(taskId, ownerId)

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] Clear redis task cache failed with error: %v\n", errRedis)
	}
	return nil
}

func (s *taskService) notifyUpdateTaskStatus(taskId uuid.UUID, ownerId uuid.UUID) {
	queueKey := "queue:notifications:status"
	cxt := context.Background()
	task, err := s.taskRepository.FindById(cxt, taskId, ownerId)

	if err != nil {
		fmt.Printf("[ERROR] publish task update status event failed: %v\n", err)
		return
	}

	event := events.NewUpdateTaskEvent(taskId, task.Description, task.Status, task.Assignee)

	err = s.redisClient.Push(cxt, queueKey, event)

	if err != nil {
		fmt.Printf("[ERROR] publish task update status event failed: %v\n", err)
	}

}

func (s *taskService) changeStatusIsValid(oldStatus string, newStatus string) bool {
	if oldStatus == newStatus {
		return false
	}

	if oldStatus == "DONE" {
		return false
	}

	if newStatus == "TODO" {
		return false
	}

	return true
}
