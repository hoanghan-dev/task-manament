package services

import (
	"context"
	"dev/task-management/internal/modules/workspace/dto"
	"dev/task-management/internal/modules/workspace/mapper"
	repositories "dev/task-management/internal/modules/workspace/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WorkspaceService interface {
	GetWorkspace(ctx context.Context, ownerId uuid.UUID) (*dto.WorkspaceResponseDTO, error)
	CreateWorkspaceDefault(ctx context.Context, ownerId uuid.UUID, fullName string) (*dto.WorkspaceResponseDTO, error)
	UpdateWorkspace(ctx context.Context, workspaceReq *dto.WorkspaceUpdateResquestDTO, ownerId uuid.UUID) error
	DeteleWorkspace(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) error
	IsWorkspaceOwnedBy(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) (bool, error)
}

type workspaceService struct {
	repo repositories.WorkspaceRepository
}

func NewWorkspaceService(repo repositories.WorkspaceRepository) WorkspaceService {
	return &workspaceService{
		repo: repo,
	}
}

func (s *workspaceService) GetWorkspace(ctx context.Context, ownerId uuid.UUID) (*dto.WorkspaceResponseDTO, error) {
	wp, err := s.repo.FindByOwnerId(ctx, ownerId)

	if err != nil {
		return nil, err // AppError from repository (NotFound or Internal)
	}

	return mapper.EntityToWorkspaceResponse(wp), nil
}
func (s *workspaceService) CreateWorkspaceDefault(ctx context.Context, ownerId uuid.UUID, fullName string) (*dto.WorkspaceResponseDTO, error) {
	workspaceReq := &dto.WorkspaceResquestDTO{
		Name:        fmt.Sprintf("%v's Workspace", fullName),
		Description: "A workspace to organize projects, tasks, and collaboration",
	}
	workspace := mapper.WorkspaceResquestToEntity(uuid.New(), workspaceReq, ownerId, time.Now())
	wpCreated, err := s.repo.Create(ctx, workspace)

	if err != nil {
		return nil, err // AppError from repository
	}

	return mapper.EntityToWorkspaceResponse(wpCreated), nil
}
func (s *workspaceService) UpdateWorkspace(
	ctx context.Context, workspaceReq *dto.WorkspaceUpdateResquestDTO, ownerId uuid.UUID) error {
	workspace := mapper.WorkspaceUpdateResquestToEntity(workspaceReq)

	err := s.repo.Update(ctx, workspace, workspace.WorkspaceId, ownerId)

	if err != nil {
		return err // AppError from repository (NotFound or Internal)
	}

	return nil
}
func (s *workspaceService) DeteleWorkspace(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) error {
	err := s.repo.Delete(ctx, wpId, ownerId)

	if err != nil {
		return err // AppError from repository (NotFound or Internal)
	}

	return nil
}

func (s *workspaceService) IsWorkspaceOwnedBy(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) (bool, error) {
	return s.repo.IsWorkspaceOwnedBy(ctx, wpId, ownerId)
}
