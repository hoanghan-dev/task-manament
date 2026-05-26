package dto

import "github.com/google/uuid"

type WorkspaceResquestDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type WorkspaceUpdateResquestDTO struct {
	WorkspaceId uuid.UUID `json:"workspace_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
}
