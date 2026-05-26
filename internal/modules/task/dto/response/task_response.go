package response

import (
	"time"

	"github.com/google/uuid"
)

type TaskResponse struct {
	Id          uuid.UUID `json:"task_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Assignee    uuid.UUID `json:"assignee_id" binding:"required"`
	Workspace   uuid.UUID `json:"workspace_id" binding:"required"`
	CreateAt    time.Time `json:"create_at" binding:"required"`
}
