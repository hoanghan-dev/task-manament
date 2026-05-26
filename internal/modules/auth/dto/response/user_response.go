package response

import (
	"dev/task-management/internal/modules/workspace/dto"
	"time"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	Id        uuid.UUID                `json:"user_id"`
	Email     string                   `json:"email"`
	FullName  string                   `json:"full_name"`
	Workspace dto.WorkspaceResponseDTO `json:"workspace"`
	CreateAt  time.Time                `json:"create_at"`
}

type UserAuthResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
