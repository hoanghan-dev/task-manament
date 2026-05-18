package services

import (
	"dev/task-management/internal/dto/request"
	"dev/task-management/internal/entities"
	"dev/task-management/internal/mapper"
	"dev/task-management/internal/repositories"
	"errors"

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

func (s *TaskService) GetAllTask() ([]*entities.Task, error) {
	tasks := s.taskRepository.FindAll()

	if len(tasks) == 0 {
		return nil, errors.New("task not found")
	}
	
	return tasks, nil
}

func (s *TaskService) GetTask(id uuid.UUID) (*entities.Task, error) {
	task, err := s.taskRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) CreateTask(t *request.TaskRequest) (*entities.Task, error) {
	id := uuid.New()
	task := mapper.TaskRequestToEntity(id, t)

	if !task.StatusIsValid() {
		return nil, errors.New("task status invalid.")
	}

	return s.taskRepository.CreateTask(task), nil
}

func (s *TaskService) UpdateTask(id uuid.UUID, t *request.TaskRequest) error {
	task := mapper.TaskRequestToEntity(id, t)

	err := s.taskRepository.UpdateTask(id, task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(id uuid.UUID) error {
	err := s.taskRepository.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}
