package request

import "github.com/google/uuid"

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Assignee    uuid.UUID `json:"assignee_id"`
	Workspace   uuid.UUID `json:"workspace_id" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Status      string    `json:"status" binding:"required"`
	Workspace   uuid.UUID `json:"workspace_id" binding:"required"`
}
  
type AssignTaskRequest struct {
	TaskId      uuid.UUID `json:"task_id" binding:"required"`
	WorkspaceId uuid.UUID `json:"workspace_id" binding:"required"`
	AssigneeId  uuid.UUID `json:"assignee_id" binding:"required"`
}
