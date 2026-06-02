package jobs

import "github.com/google/uuid"

type AssignTaskJob struct {
	SenderId   uuid.UUID `json:"sender_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	TaskId     uuid.UUID `json:"task_id"`
}

func NewAssignTaskJob(senderId uuid.UUID, receiverId uuid.UUID, taskId uuid.UUID) *AssignTaskJob {
	return &AssignTaskJob{
		SenderId:   senderId,
		ReceiverId: receiverId,
		TaskId:     taskId,
	}
}
