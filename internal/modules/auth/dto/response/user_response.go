package response

import (
	"time"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	Id       uuid.UUID `json:"user_id" binding:"required, email"`
	Email    string    `json:"emai" binding:"required, email"`
	FullName string    `json:"full_name" binding:"required"`
	CreateAt time.Time `json:"create_at" binding:"required"`
}

type UserAuthResponseDTO struct {
	AccessToken  string `json:"access_token" binding:"required"`
	RefreshToken string `json:"refresh_token" binding:"required"`
}
