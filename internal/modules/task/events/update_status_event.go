package events

import "github.com/google/uuid"

type UpdateTaskEvent struct {
	TaskId      uuid.UUID `json:"task_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ReceiverId  uuid.UUID `json:"receiver_id"`
}

func NewUpdateTaskEvent(TaskId uuid.UUID, Description string, Status string, ReceiverId uuid.UUID) *UpdateTaskEvent {
	return &UpdateTaskEvent{
		TaskId:      TaskId,
		Description: Description,
		Status:      Status,
		ReceiverId:  ReceiverId,
	}
}
