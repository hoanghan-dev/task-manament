package services

import (
	"context"
	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/dto/response"
	"dev/task-management/internal/modules/task/mapper"
	"dev/task-management/internal/modules/task/repositories"
	"dev/task-management/internal/modules/workspace/services"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
	GetAllTask(ctx context.Context, ownerId uuid.UUID) ([]response.TaskResponse, error)
	GetTask(context context.Context, id uuid.UUID, ownerId uuid.UUID) (*response.TaskResponse, error)
	CreateTask(ctx context.Context, t *request.TaskRequest, ownerId uuid.UUID) (*response.TaskResponse, error)
	UpdateTask(ctx context.Context, id uuid.UUID, t *request.TaskRequest, ownerId uuid.UUID) error
	DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error
}

type taskService struct {
	taskRepository   repositories.TaskRepository
	workspaceService services.WorkspaceService
}

func NewTaskService(repo repositories.TaskRepository, workspaceService services.WorkspaceService) TaskService {
	return &taskService{
		taskRepository:   repo,
		workspaceService: workspaceService,
	}
}

func (s *taskService) GetAllTask(ctx context.Context, ownerId uuid.UUID) ([]response.TaskResponse, error) {
	tasks, err := s.taskRepository.FindAll(ctx, ownerId)

	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, nil
	}

	taskRes := make([]response.TaskResponse, 0)

	for _, task := range tasks {
		t := mapper.EntityToTaskResponse(&task)
		taskRes = append(taskRes, *t)
	}
	return taskRes, nil
}

func (s *taskService) GetTask(context context.Context, id uuid.UUID, ownerId uuid.UUID) (*response.TaskResponse, error) {
	task, err := s.taskRepository.FindById(context, id, ownerId)
	if err != nil {
		return nil, err
	}
	taskRes := mapper.EntityToTaskResponse(task)
	return taskRes, nil
}

func (s *taskService) CreateTask(ctx context.Context, t *request.TaskRequest, ownerId uuid.UUID) (*response.TaskResponse, error) {

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

	return mapper.EntityToTaskResponse(task), nil
}

func (s *taskService) UpdateTask(ctx context.Context, id uuid.UUID, t *request.TaskRequest, ownerId uuid.UUID) error {
	task := mapper.TaskUpdateToEntity(t)
	err := s.taskRepository.UpdateTask(ctx, id, task, ownerId)
	if err != nil {
		return err
	}
	return nil
}

func (s *taskService) DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	err := s.taskRepository.DeleteTask(ctx, id, ownerId)
	if err != nil {
		return err
	}
	return nil
}
