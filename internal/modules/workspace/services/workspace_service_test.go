package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"dev/task-management/internal/modules/workspace/dto"
	entities "dev/task-management/internal/modules/workspace/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeWorkspaceRepository struct {
	findByOwnerIdCalled bool
	findByOwnerIdArg    uuid.UUID
	findByOwnerIdResult *entities.Workspace
	findByOwnerIdErr    error

	createCalled bool
	createArg    *entities.Workspace
	createResult *entities.Workspace
	createErr    error

	updateCalled bool
	updateArgWP  *entities.Workspace
	updateArgID  uuid.UUID
	updateOwner  uuid.UUID
	updateErr    error

	deleteCalled bool
	deleteArgID  uuid.UUID
	deleteOwner  uuid.UUID
	deleteErr    error

	isOwnedCalled bool
	isOwnedWPID   uuid.UUID
	isOwnedOwner  uuid.UUID
	isOwnedResult bool
	isOwnedErr    error
}

func (f *fakeWorkspaceRepository) FindByOwnerId(ctx context.Context, ownerId uuid.UUID) (*entities.Workspace, error) {
	f.findByOwnerIdCalled = true
	f.findByOwnerIdArg = ownerId
	return f.findByOwnerIdResult, f.findByOwnerIdErr
}

func (f *fakeWorkspaceRepository) Create(ctx context.Context, wp *entities.Workspace) (*entities.Workspace, error) {
	f.createCalled = true
	f.createArg = wp

	if f.createErr != nil {
		return nil, f.createErr
	}

	if f.createResult != nil {
		return f.createResult, nil
	}

	return wp, nil
}

func (f *fakeWorkspaceRepository) Update(ctx context.Context, wp *entities.Workspace, wpId uuid.UUID, ownerId uuid.UUID) error {
	f.updateCalled = true
	f.updateArgWP = wp
	f.updateArgID = wpId
	f.updateOwner = ownerId
	return f.updateErr
}

func (f *fakeWorkspaceRepository) Delete(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) error {
	f.deleteCalled = true
	f.deleteArgID = wpId
	f.deleteOwner = ownerId
	return f.deleteErr
}

func (f *fakeWorkspaceRepository) IsWorkspaceOwnedBy(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) (bool, error) {
	f.isOwnedCalled = true
	f.isOwnedWPID = wpId
	f.isOwnedOwner = ownerId
	return f.isOwnedResult, f.isOwnedErr
}

func TestGetWorkspace_Success(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()

	repo := &fakeWorkspaceRepository{
		findByOwnerIdResult: &entities.Workspace{
			WorkspaceId: workspaceId,
			Name:        "Main Workspace",
			Description: "Workspace description",
			OwnerId:     ownerId,
			CreateAt:    time.Now(),
		},
	}

	service := NewWorkspaceService(repo)

	result, err := service.GetWorkspace(ctx, ownerId)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, repo.findByOwnerIdCalled)
	assert.Equal(t, ownerId, repo.findByOwnerIdArg)

	assert.Equal(t, workspaceId, result.WorkspaceId)
	assert.Equal(t, "Main Workspace", result.Name)
	assert.Equal(t, "Workspace description", result.Description)
}

func TestGetWorkspace_RepositoryError(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	expectedErr := errors.New("repository error")

	repo := &fakeWorkspaceRepository{
		findByOwnerIdErr: expectedErr,
	}

	service := NewWorkspaceService(repo)

	result, err := service.GetWorkspace(ctx, ownerId)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, repo.findByOwnerIdCalled)
	assert.Equal(t, ownerId, repo.findByOwnerIdArg)
}

func TestCreateWorkspaceDefault_Success(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	fullName := "Giyuu"

	repo := &fakeWorkspaceRepository{}

	service := NewWorkspaceService(repo)

	result, err := service.CreateWorkspaceDefault(ctx, ownerId, fullName)

	require.NoError(t, err)
	require.NotNil(t, result)

	require.True(t, repo.createCalled)
	require.NotNil(t, repo.createArg)

	assert.NotEqual(t, uuid.Nil, repo.createArg.WorkspaceId)
	assert.Equal(t, ownerId, repo.createArg.OwnerId)
	assert.Equal(t, "Giyuu's Workspace", repo.createArg.Name)
	assert.Equal(t, "A workspace to organize projects, tasks, and collaboration", repo.createArg.Description)
	assert.False(t, repo.createArg.CreateAt.IsZero())

	assert.Equal(t, repo.createArg.WorkspaceId, result.WorkspaceId)
	assert.Equal(t, repo.createArg.Name, result.Name)
	assert.Equal(t, repo.createArg.Description, result.Description)
}

func TestCreateWorkspaceDefault_RepositoryError(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	expectedErr := errors.New("create workspace failed")

	repo := &fakeWorkspaceRepository{
		createErr: expectedErr,
	}

	service := NewWorkspaceService(repo)

	result, err := service.CreateWorkspaceDefault(ctx, ownerId, "Giyuu")

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, expectedErr)

	assert.True(t, repo.createCalled)
	require.NotNil(t, repo.createArg)
	assert.Equal(t, ownerId, repo.createArg.OwnerId)
}

func TestUpdateWorkspace_Success(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()

	repo := &fakeWorkspaceRepository{}
	service := NewWorkspaceService(repo)

	req := &dto.WorkspaceUpdateResquestDTO{
		WorkspaceId: workspaceId,
		Name:        "Updated Workspace",
		Description: "Updated description",
	}

	err := service.UpdateWorkspace(ctx, req, ownerId)

	require.NoError(t, err)

	require.True(t, repo.updateCalled)
	require.NotNil(t, repo.updateArgWP)

	assert.Equal(t, workspaceId, repo.updateArgID)
	assert.Equal(t, ownerId, repo.updateOwner)
	assert.Equal(t, workspaceId, repo.updateArgWP.WorkspaceId)
	assert.Equal(t, "Updated Workspace", repo.updateArgWP.Name)
	assert.Equal(t, "Updated description", repo.updateArgWP.Description)
}

func TestUpdateWorkspace_RepositoryError(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()
	expectedErr := errors.New("update workspace failed")

	repo := &fakeWorkspaceRepository{
		updateErr: expectedErr,
	}
	service := NewWorkspaceService(repo)

	req := &dto.WorkspaceUpdateResquestDTO{
		WorkspaceId: workspaceId,
		Name:        "Updated Workspace",
		Description: "Updated description",
	}

	err := service.UpdateWorkspace(ctx, req, ownerId)

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)

	assert.True(t, repo.updateCalled)
	assert.Equal(t, workspaceId, repo.updateArgID)
	assert.Equal(t, ownerId, repo.updateOwner)
}

func TestDeteleWorkspace_Success(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()

	repo := &fakeWorkspaceRepository{}
	service := NewWorkspaceService(repo)

	err := service.DeteleWorkspace(ctx, workspaceId, ownerId)

	require.NoError(t, err)

	assert.True(t, repo.deleteCalled)
	assert.Equal(t, workspaceId, repo.deleteArgID)
	assert.Equal(t, ownerId, repo.deleteOwner)
}

func TestDeteleWorkspace_RepositoryError(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()
	expectedErr := errors.New("delete workspace failed")

	repo := &fakeWorkspaceRepository{
		deleteErr: expectedErr,
	}
	service := NewWorkspaceService(repo)

	err := service.DeteleWorkspace(ctx, workspaceId, ownerId)

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)

	assert.True(t, repo.deleteCalled)
	assert.Equal(t, workspaceId, repo.deleteArgID)
	assert.Equal(t, ownerId, repo.deleteOwner)
}

func TestIsWorkspaceOwnedBy_ReturnsTrue(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()

	repo := &fakeWorkspaceRepository{
		isOwnedResult: true,
	}
	service := NewWorkspaceService(repo)

	owned, err := service.IsWorkspaceOwnedBy(ctx, workspaceId, ownerId)

	require.NoError(t, err)
	assert.True(t, owned)

	assert.True(t, repo.isOwnedCalled)
	assert.Equal(t, workspaceId, repo.isOwnedWPID)
	assert.Equal(t, ownerId, repo.isOwnedOwner)
}

func TestIsWorkspaceOwnedBy_ReturnsFalse(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()

	repo := &fakeWorkspaceRepository{
		isOwnedResult: false,
	}
	service := NewWorkspaceService(repo)

	owned, err := service.IsWorkspaceOwnedBy(ctx, workspaceId, ownerId)

	require.NoError(t, err)
	assert.False(t, owned)

	assert.True(t, repo.isOwnedCalled)
	assert.Equal(t, workspaceId, repo.isOwnedWPID)
	assert.Equal(t, ownerId, repo.isOwnedOwner)
}

func TestIsWorkspaceOwnedBy_RepositoryError(t *testing.T) {
	ctx := context.Background()
	ownerId := uuid.New()
	workspaceId := uuid.New()
	expectedErr := errors.New("ownership check failed")

	repo := &fakeWorkspaceRepository{
		isOwnedErr: expectedErr,
	}
	service := NewWorkspaceService(repo)

	owned, err := service.IsWorkspaceOwnedBy(ctx, workspaceId, ownerId)

	require.Error(t, err)
	assert.False(t, owned)
	assert.ErrorIs(t, err, expectedErr)

	assert.True(t, repo.isOwnedCalled)
	assert.Equal(t, workspaceId, repo.isOwnedWPID)
	assert.Equal(t, ownerId, repo.isOwnedOwner)
}
