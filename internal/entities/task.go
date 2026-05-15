package entities

import "github.com/google/uuid"

type Task struct {
	Id          uuid.UUID `json:"id" validate:"required"`
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	Status      string    `json:"status" validate:"required"`
	Assignee    string    `json:"assignee" validate:"required"`
}

func NewTask(id uuid.UUID, title string, description string, status string, assignee string) *Task {
	return &Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      status,
		Assignee:    assignee,
	}
}

var statusArr = [3]string{"TODO", "IN_PROGRESS", "DONE"}

func (t *Task) StatusIsValid() bool {
	for _, s := range statusArr {
		if s == t.Status {
			return true
		}
	}
	return false
}
