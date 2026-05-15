package repositories

import (
	"dev/task-management/internal/db"
	"dev/task-management/internal/entities"
	"errors"

	"github.com/google/uuid"
)

type TaskRepository interface {
	FindAll() []*entities.Task
	FindById(id uuid.UUID) (*entities.Task, error)
	CreateTask(t *entities.Task) *entities.Task
	UpdateTask(id uuid.UUID, task *entities.Task) error
	DeleteTask(id uuid.UUID) error
}

type taskRepository struct {
}

func NewTaskRepository() TaskRepository {
	return &taskRepository{}
}

func (r *taskRepository) FindAll() []*entities.Task {
	return db.TaskDataMock
}

func (r *taskRepository) FindById(id uuid.UUID) (*entities.Task, error) {

	for _, task := range db.TaskDataMock {
		if task.Id == id {
			return task, nil
		}
	}
	return nil, errors.New("task not found.")
}

func (r *taskRepository) CreateTask(t *entities.Task) *entities.Task {
	db.TaskDataMock = append(db.TaskDataMock, t)
	return t
}

func (r *taskRepository) UpdateTask(id uuid.UUID, t *entities.Task) error {
	for i, task := range db.TaskDataMock {
		if task.Id == id {
			db.TaskDataMock[i] = t
			return nil
		}
	}
	return errors.New("task not found to updated.")
}

func (r *taskRepository) DeleteTask(id uuid.UUID) error {
	for i, task := range db.TaskDataMock {
		if task.Id == id {
			//slice first (điểm trước element cần remove) = db.TaskDataMock[:i] \
			//slice end (điểm sau element cần remove) = db.TaskDataMock[i+1:]
			// ... là tất cả các element phía sau
			db.TaskDataMock = append(db.TaskDataMock[:i], db.TaskDataMock[i+1:]...)
			return nil
		}
	}
	return errors.New("task not found to delete.")
}
