package mapper

import (
	dto "dev/task-management/internal/modules/workspace/dto"
	entities "dev/task-management/internal/modules/workspace/entities"
	"time"

	"github.com/google/uuid"
)

func WorkspaceResquestToEntity(id uuid.UUID, wq *dto.WorkspaceResquestDTO, ownerId uuid.UUID, createAt time.Time) *entities.Workspace {
	return &entities.Workspace{
		WorkspaceId: id,
		Name:        wq.Name,
		Description: wq.Description,
		OwnerId:     ownerId,
		CreateAt:    createAt,
	}
}

func WorkspaceUpdateResquestToEntity(wq *dto.WorkspaceUpdateResquestDTO) *entities.Workspace {
	return &entities.Workspace{
		WorkspaceId: wq.WorkspaceId,
		Name:        wq.Name,
		Description: wq.Description,
	}
}

func EntityToWorkspaceResponse(wp *entities.Workspace) *dto.WorkspaceResponseDTO {
	return &dto.WorkspaceResponseDTO{
		WorkspaceId: wp.WorkspaceId,
		Name:        wp.Name,
		Description: wp.Description,
		CreateAt:    wp.CreateAt,
	}
}
