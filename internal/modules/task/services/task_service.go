package services

import (
	"context"
	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/dto/response"
	"dev/task-management/internal/modules/task/mapper"
	"dev/task-management/internal/modules/task/repositories"
	"errors"
	"time"

	"github.com/google/uuid"
)

type TaskService struct {
	taskRepository repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: repo,
	}
}

func (s *TaskService) GetAllTask(ctx context.Context) ([]response.TaskResponse, error) {
	tasks, err := s.taskRepository.FindAll(ctx)

	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, errors.New("task not found")
	}

	taskRes := make([]response.TaskResponse, 0)

	for _, task := range tasks {
		t := mapper.EntityToTaskResponse(&task)
		taskRes = append(taskRes, *t)
	}
	return taskRes, nil
}

func (s *TaskService) GetTask(context context.Context, id uuid.UUID) (*response.TaskResponse, error) {
	task, err := s.taskRepository.FindById(context, id)
	if err != nil {
		return nil, err
	}
	taskRes := mapper.EntityToTaskResponse(task)
	return taskRes, nil
}

func (s *TaskService) CreateTask(ctx context.Context, t *request.TaskRequest) (*response.TaskResponse, error) {
	id := uuid.New()
	createAt := time.Now()
	task := mapper.TaskRequestToEntity(id, t, createAt)

	if !task.StatusIsValid() {
		return nil, errors.New("task status invalid.")
	}

	err := s.taskRepository.CreateTask(ctx, task)

	if err != nil {
		return nil, errors.New("Task creation failed.")
	}

	return mapper.EntityToTaskResponse(task), nil
}

func (s *TaskService) UpdateTask(ctx context.Context, id uuid.UUID, t *request.TaskRequest) error {
	task := mapper.TaskUpdateToEntity(t)
	err := s.taskRepository.UpdateTask(ctx, id, task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(ctx context.Context, id uuid.UUID) error {
	err := s.taskRepository.DeleteTask(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
