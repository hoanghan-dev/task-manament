package repositories

import (
	"dev/task-management/internal/entities"
	"errors"
)

type TaskRepository interface {
	FindAll() []*entities.Task
	FindById(id int) (*entities.Task, error)
	CreateTask(t entities.Task) *entities.Task
	UpdateTask(id int, task *entities.Task) (*entities.Task, error)
	DeleteTask(id int) error
}

type taskRepository struct {
	tasks []*entities.Task
}

func NewtaskRepository() TaskRepository {
	r := &taskRepository{}
	r.initTask()
	return r
}

func (r *taskRepository) initTask() {
	t1 := entities.NewTask(1, "Học Go Lang", "Học syntax kỹ zô..")
	t2 := entities.NewTask(2, "Học Go Http", "Học api kỹ zô..")
	t3 := entities.NewTask(3, "Học Go Security", "Học jwt kỹ zô..")
	r.tasks = append(r.tasks, t1, t2, t3)
}

func (r *taskRepository) FindAll() []*entities.Task {
	return r.tasks
}

func (r *taskRepository) FindById(id int) (*entities.Task, error) {
	for _, task := range r.tasks {
		if task.Id == id {
			return task, nil
		}
	}
	return nil, errors.New("Task not found!")
}

func (r *taskRepository) CreateTask(t entities.Task) *entities.Task {
	r.tasks = append(r.tasks, &t)
	return &t
}

func (r *taskRepository) UpdateTask(id int, t *entities.Task) (*entities.Task, error) {
	task, err := r.FindById(id)

	if err != nil {
		return nil, err
	}

	task.Title = t.Title
	task.Description = t.Description

	return task, nil
}

func (r *taskRepository) DeleteTask(id int) error {
	for i, task := range r.tasks {
		if task.Id == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("deleted faid!")
}
