package dto

import (
	"time"

	"github.com/google/uuid"
)

type WorkspaceResponseDTO struct {
	WorkspaceId uuid.UUID `json:"workspace_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"create_at"`
}
