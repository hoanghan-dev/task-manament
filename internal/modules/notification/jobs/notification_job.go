package jobs

import "github.com/google/uuid"

type NotificationJob struct {
	SenderId   uuid.UUID `json:"sender_id"`
	ReceiverId uuid.UUID `json:"receiver_id"`
	TaskId     uuid.UUID `json:"task_id"`
}

func NewNotificationJob(senderId uuid.UUID, receiverId uuid.UUID, taskId uuid.UUID) *NotificationJob {
	return &NotificationJob{
		SenderId:   senderId,
		ReceiverId: receiverId,
		TaskId:     taskId,
	}
}
