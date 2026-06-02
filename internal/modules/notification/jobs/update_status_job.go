package jobs

import "github.com/google/uuid"

type UpdateStatusJob struct {
	TaskId      uuid.UUID `json:"task_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	ReceiverId  uuid.UUID `json:"receiver_id"`
}

func NewUpdateStatusJob(TaskId uuid.UUID, Description string, Status string, ReceiverId uuid.UUID) *UpdateStatusJob {
	return &UpdateStatusJob{
		TaskId:      TaskId,
		Description: Description,
		Status:      Status,
		ReceiverId:  ReceiverId,
	}
}
