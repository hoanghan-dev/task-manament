package entities

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	Id          uuid.UUID
	Title       string
	Description string
	Status      string
	Assignee    uuid.UUID
	Workspace   uuid.UUID
	CreateAt    time.Time
}

func  NewTask(id uuid.UUID, title string, description string, status string, assignee uuid.UUID, workspace uuid.UUID, createAt time.Time) *Task {
	return &Task{
		Id:          id,
		Title:       title,
		Description: description,
		Status:      status,
		Assignee:    assignee,
		Workspace:   workspace,
		CreateAt:    createAt,
	}
}

var statusArr = [4]string{"TODO", "IN_PROGRESS", "DONE", "BLOCKED"}

func (t *Task) StatusIsValid() bool {
	for _, s := range statusArr {
		if s == t.Status {
			return true
		}
	}
	return false
}
