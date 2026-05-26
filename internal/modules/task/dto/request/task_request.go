package request

import "github.com/google/uuid"

type TaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Assignee    uuid.UUID `json:"assignee_id"`
	Workspace   uuid.UUID `json:"workspace_id" binding:"required"`
}
