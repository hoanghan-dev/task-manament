package services

import (
	"context"
	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/dto/response"
	"dev/task-management/internal/modules/task/events"
	"dev/task-management/internal/modules/task/mapper"
	"dev/task-management/internal/modules/task/repositories"
	workspaceService "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/pkg/cache"
	"errors"
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
		fmt.Printf("[ERROR] Failed to get data in redis with error: %v\n", err)
		fmt.Println("[INFO] Jump into cache")
		return tasksCache, nil
	}
	fmt.Println("[INFO] Get data into database")
	tasks, err := s.taskRepository.FindAll(ctx, ownerId)

	if err != nil {
		return nil, err
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

func (s *taskService) GetTask(context context.Context, id uuid.UUID, ownerId uuid.UUID) (*response.TaskResponse, error) {

	redisKey := "tasks:owner_id:" + ownerId.String() + ":task_id:" + id.String()

	taskCache := &response.TaskResponse{}

	err := s.redisClient.Get(redisKey, &taskCache)

	if err == nil {
		return taskCache, nil
	}

	task, err := s.taskRepository.FindById(context, id, ownerId)
	if err != nil {
		fmt.Printf("[ERROR] Set task into cache failed with error: %v\n", err)
	}
	taskRes := mapper.EntityToTaskResponse(task)

	s.redisClient.Set(redisKey, taskRes, 10)
	return taskRes, nil
}

func (s *taskService) CreateTask(ctx context.Context, t *request.CreateTaskRequest, ownerId uuid.UUID) (*response.TaskResponse, error) {

	ownered, err := s.workspaceService.IsWorkspaceOwnedBy(ctx, t.Workspace, ownerId)

	if err != nil {
		return nil, errors.New("workspace not found")
	}

	if !ownered {
		return nil, errors.New("access denied")
	}

	id := uuid.New()
	createAt := time.Now()
	task := mapper.TaskRequestToEntity(id, t, createAt)

	if !task.StatusIsValid() {
		return nil, errors.New("task status invalid.")
	}

	createErr := s.taskRepository.CreateTask(ctx, task)

	if createErr != nil {
		return nil, errors.New("Task creation failed.")
	}

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] CLear redis task cache failed with error: %v\n", errRedis)
	}

	return mapper.EntityToTaskResponse(task), nil
}

func (s *taskService) UpdateTask(ctx context.Context, id uuid.UUID, t *request.UpdateTaskRequest, ownerId uuid.UUID) error {
	task := mapper.TaskUpdateToEntity(t)
	err := s.taskRepository.UpdateTask(ctx, id, task, ownerId)
	if err != nil {
		return err
	}

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] CLear redis task cache failed with error: %v\n", errRedis)
	}
	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	err := s.taskRepository.DeleteTask(ctx, id, ownerId)
	if err != nil {
		return err
	}

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] CLear redis task cache failed with error: %v\n", errRedis)
	}
	return nil
}

func (s *taskService) AssignTask(ctx context.Context, assignTaskReq *request.AssignTaskRequest, ownerId uuid.UUID) error {
	if !s.taskRepository.TaskExistsInWorkspace(ctx, assignTaskReq.TaskId, assignTaskReq.WorkspaceId) {
		return errors.New("task id not exists in workspace")
	}

	exists := s.taskRepository.TaskIsExists(ctx, assignTaskReq.TaskId)

	if !exists {
		return fmt.Errorf("Task not found with id: %v", assignTaskReq.TaskId)
	}

	ownered, err := s.workspaceService.IsWorkspaceOwnedBy(ctx, assignTaskReq.WorkspaceId, ownerId)

	if err != nil {
		return err
	}

	if !ownered {
		return errors.New("access denied")
	}

	errAssign := s.taskRepository.AssignTask(ctx, assignTaskReq.TaskId, assignTaskReq.AssigneeId)

	if errAssign != nil {
		return errAssign
	}

	go s.notifyAssignTask(ownerId, assignTaskReq.AssigneeId, assignTaskReq.TaskId)

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] CLear redis task cache failed with error: %v\n", errRedis)
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
		return errors.New("task id not exists in workspace")
	}

	task, err := s.taskRepository.FindById(ctx, taskId, ownerId)

	if err != nil {
		return fmt.Errorf("Task not found with id: %v", taskId)
	}

	ownered, err := s.workspaceService.IsWorkspaceOwnedBy(ctx, taskReq.WorkspaceId, ownerId)

	if err != nil {
		return err
	}

	if !ownered {
		return errors.New("access denied")
	}

	if !s.changeStatusIsValid(task.Status, taskReq.Status) {
		return fmt.Errorf("can not change task status from %s to %s", task.Status, taskReq.Status)
	}

	err = s.taskRepository.UpdateTaskStatus(ctx, taskId, taskReq.Status)

	if err != nil {
		return err
	}

	s.notifyUpdateTaskStatus(taskId, ownerId)

	errRedis := s.redisClient.Clear("tasks:*")

	if errRedis != nil {
		fmt.Printf("[ERROR] CLear redis task cache failed with error: %v\n", errRedis)
	}
	return nil
}

func (s *taskService) notifyUpdateTaskStatus(taskId uuid.UUID, ownerId uuid.UUID) {
	queueKey := "queue:notifications:status"
	cxt := context.Background()
	task, err := s.taskRepository.FindById(cxt, taskId, ownerId)

	if err != nil {
		fmt.Printf("[ERROR] publish task update status event failed: %v\n", err)
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
