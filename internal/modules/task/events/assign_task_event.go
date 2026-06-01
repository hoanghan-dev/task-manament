package events

import "github.com/google/uuid"

type AssignTaskEvent struct {
	SenderId   uuid.UUID `json:"sender_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	TaskId     uuid.UUID `json:"task_id"`
}

func NewAssignTaskEvent(senderId uuid.UUID, receiverId uuid.UUID, taskId uuid.UUID) *AssignTaskEvent {
	return &AssignTaskEvent{
		SenderId:   senderId,
		ReceiverId: receiverId,
		TaskId:     taskId,
	}
}
