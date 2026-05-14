package services

import (
	"dev/task-management/internal/entities"
	"dev/task-management/internal/repositories"
)

type TaskService struct {
	taskRepository repositories.TaskRepository
}

func NewTaskService(repo repositories.TaskRepository) *TaskService {
	return &TaskService{
		taskRepository: repo,
	}
}

func (s *TaskService) GetAllTask() []*entities.Task {
	return s.taskRepository.FindAll()
}

func (s *TaskService) GetTask(id int) (*entities.Task, error) {
	task, err := s.taskRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) CreateTask(t entities.Task) *entities.Task {
	return s.taskRepository.CreateTask(t)
}

func (s *TaskService) UpdateTask(id int, t *entities.Task) (*entities.Task, error) {
	task, err := s.taskRepository.UpdateTask(id, t)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *TaskService) DeleteTask(id int) error {
	err := s.taskRepository.DeleteTask(id)
	if err != nil {
		return err
	}
	return nil
}
