package mapper

import (
	"dev/task-management/internal/modules/auth/dto/request"
	"dev/task-management/internal/modules/auth/dto/response"
	"dev/task-management/internal/modules/auth/entities"
	"dev/task-management/internal/modules/workspace/dto"
	"time"

	"github.com/google/uuid"
)

func RegisterUserRequestToEntity(req *request.RegisterUserRequestDTO) *entities.User {

	return &entities.User{
		Id:       uuid.New(),
		Email:    req.Email,
		Password: req.Password,
		FullName: req.FullName,
		CreateAt: time.Now(),
	}
}

func EntityToUserResponse(u *entities.User, wp *dto.WorkspaceResponseDTO) *response.UserResponseDTO {
	return &response.UserResponseDTO{
		Id:        u.Id,
		Email:     u.Email,
		FullName:  u.FullName,
		Workspace: *wp,
		CreateAt:  u.CreateAt,
	}
}
