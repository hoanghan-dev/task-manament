package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"dev/task-management/internal/modules/auth/dto/request"
	authEntities "dev/task-management/internal/modules/auth/entities"
	workspaceDTO "dev/task-management/internal/modules/workspace/dto"
	"dev/task-management/pkg/apperror"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type fakeUserRepository struct {
	getUserByEmailCalled bool
	getUserByEmailArg    string
	getUserByEmailResult *authEntities.User
	getUserByEmailErr    error

	createUserCalled bool
	createUserArg    *authEntities.User
	createUserErr    error

	userIsExistsCalled bool
	userIsExistsArg    uuid.UUID
	userIsExistsResult bool
}

func (f *fakeUserRepository) GetUserByEmail(ctx context.Context, email string) (*authEntities.User, error) {
	f.getUserByEmailCalled = true
	f.getUserByEmailArg = email
	return f.getUserByEmailResult, f.getUserByEmailErr
}

func (f *fakeUserRepository) CreateUser(ctx context.Context, user *authEntities.User) error {
	f.createUserCalled = true
	f.createUserArg = user
	return f.createUserErr
}

func (f *fakeUserRepository) UserIsExists(ctx context.Context, userId uuid.UUID) bool {
	f.userIsExistsCalled = true
	f.userIsExistsArg = userId
	return f.userIsExistsResult
}

type fakeWorkspaceService struct {
	getWorkspaceCalled bool
	getWorkspaceResult *workspaceDTO.WorkspaceResponseDTO
	getWorkspaceErr    error

	createDefaultCalled  bool
	createDefaultOwnerID uuid.UUID
	createDefaultName    string
	createDefaultResult  *workspaceDTO.WorkspaceResponseDTO
	createDefaultErr     error

	updateWorkspaceCalled bool
	updateWorkspaceErr    error

	deleteWorkspaceCalled bool
	deleteWorkspaceErr    error

	isOwnedCalled bool
	isOwnedResult bool
	isOwnedErr    error
}

func (f *fakeWorkspaceService) GetWorkspace(ctx context.Context, ownerId uuid.UUID) (*workspaceDTO.WorkspaceResponseDTO, error) {
	f.getWorkspaceCalled = true
	return f.getWorkspaceResult, f.getWorkspaceErr
}

func (f *fakeWorkspaceService) CreateWorkspaceDefault(ctx context.Context, ownerId uuid.UUID, fullName string) (*workspaceDTO.WorkspaceResponseDTO, error) {
	f.createDefaultCalled = true
	f.createDefaultOwnerID = ownerId
	f.createDefaultName = fullName

	if f.createDefaultErr != nil {
		return nil, f.createDefaultErr
	}

	if f.createDefaultResult != nil {
		return f.createDefaultResult, nil
	}

	return &workspaceDTO.WorkspaceResponseDTO{}, nil
}

func (f *fakeWorkspaceService) UpdateWorkspace(ctx context.Context, workspaceReq *workspaceDTO.WorkspaceUpdateResquestDTO, ownerId uuid.UUID) error {
	f.updateWorkspaceCalled = true
	return f.updateWorkspaceErr
}

func (f *fakeWorkspaceService) DeteleWorkspace(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) error {
	f.deleteWorkspaceCalled = true
	return f.deleteWorkspaceErr
}

func (f *fakeWorkspaceService) IsWorkspaceOwnedBy(ctx context.Context, wpId uuid.UUID, ownerId uuid.UUID) (bool, error) {
	f.isOwnedCalled = true
	return f.isOwnedResult, f.isOwnedErr
}

func validRegisterRequest() *request.RegisterUserRequestDTO {
	return &request.RegisterUserRequestDTO{
		Email:    "giyuu@example.com",
		Password: "Password123!",
		FullName: "Giyuu Tomioka",
	}
}

func validLoginRequest(password string) *request.LoginUserRequestDTO {
	return &request.LoginUserRequestDTO{
		Email:    "giyuu@example.com",
		Password: password,
	}
}

func TestRegister_Success(t *testing.T) {
	ctx := context.Background()

	repo := &fakeUserRepository{
		getUserByEmailErr: apperror.NewNotFound("user"),
	}

	workspace := &fakeWorkspaceService{}

	service := NewAuthService(repo, workspace)

	result, err := service.Register(ctx, validRegisterRequest())

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, repo.getUserByEmailCalled)
	assert.Equal(t, "giyuu@example.com", repo.getUserByEmailArg)

	require.True(t, repo.createUserCalled)
	require.NotNil(t, repo.createUserArg)
	assert.Equal(t, "giyuu@example.com", repo.createUserArg.Email)
	assert.Equal(t, "Giyuu Tomioka", repo.createUserArg.FullName)
	assert.NotEqual(t, "Password123!", repo.createUserArg.Password, "password must be hashed before saving")

	err = bcrypt.CompareHashAndPassword([]byte(repo.createUserArg.Password), []byte("Password123!"))
	assert.NoError(t, err)

	assert.True(t, workspace.createDefaultCalled)
	assert.Equal(t, repo.createUserArg.Id, workspace.createDefaultOwnerID)
	assert.Equal(t, repo.createUserArg.FullName, workspace.createDefaultName)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	ctx := context.Background()

	repo := &fakeUserRepository{
		getUserByEmailResult: &authEntities.User{
			Id:       uuid.New(),
			Email:    "giyuu@example.com",
			Password: "hashed-password",
			FullName: "Giyuu Tomioka",
			CreateAt: time.Now(),
		},
		getUserByEmailErr: nil,
	}

	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	result, err := service.Register(ctx, validRegisterRequest())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.getUserByEmailCalled)
	assert.False(t, repo.createUserCalled)
	assert.False(t, workspace.createDefaultCalled)
}

func TestRegister_InvalidPassword(t *testing.T) {
	ctx := context.Background()

	repo := &fakeUserRepository{
		getUserByEmailErr: apperror.NewNotFound("user"),
	}
	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	req := validRegisterRequest()
	req.Password = "123"

	result, err := service.Register(ctx, req)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.getUserByEmailCalled)
	assert.False(t, repo.createUserCalled)
	assert.False(t, workspace.createDefaultCalled)
}

func TestRegister_CreateUserRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("create user failed")

	repo := &fakeUserRepository{
		getUserByEmailErr: apperror.NewNotFound("user"),
		createUserErr:     expectedErr,
	}
	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	result, err := service.Register(ctx, validRegisterRequest())

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.True(t, repo.createUserCalled)
	assert.False(t, workspace.createDefaultCalled)
}

func TestRegister_CreateWorkspaceDefaultError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("create workspace failed")

	repo := &fakeUserRepository{
		getUserByEmailErr: apperror.NewNotFound("user"),
	}
	workspace := &fakeWorkspaceService{
		createDefaultErr: expectedErr,
	}
	service := NewAuthService(repo, workspace)

	result, err := service.Register(ctx, validRegisterRequest())

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.True(t, repo.createUserCalled)
	assert.True(t, workspace.createDefaultCalled)
}

func TestLogin_Success(t *testing.T) {
	ctx := context.Background()

	// Set common JWT env names used by many projects.
	// If your utils.GenerateAccessToken uses another env name, add it here.
	t.Setenv("JWT_SECRET", "unit-test-secret")
	t.Setenv("JWT_ACCESS_SECRET", "unit-test-secret")
	t.Setenv("ACCESS_TOKEN_SECRET", "unit-test-secret")

	password := "Password123!"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	require.NoError(t, err)

	repo := &fakeUserRepository{
		getUserByEmailResult: &authEntities.User{
			Id:       uuid.New(),
			Email:    "giyuu@example.com",
			Password: string(hash),
			FullName: "Giyuu Tomioka",
			CreateAt: time.Now(),
		},
	}

	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	result, err := service.Login(ctx, validLoginRequest(password))

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.NotEmpty(t, result.AccessToken)
	assert.NotEmpty(t, result.RefreshToken)
	assert.True(t, repo.getUserByEmailCalled)
	assert.Equal(t, "giyuu@example.com", repo.getUserByEmailArg)
}

func TestLogin_UserNotFound_ReturnUnauthorized(t *testing.T) {
	ctx := context.Background()

	repo := &fakeUserRepository{
		getUserByEmailErr: apperror.NewNotFound("user"),
	}

	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	result, err := service.Login(ctx, validLoginRequest("Password123!"))

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.getUserByEmailCalled)
}

func TestLogin_UserRepositoryError_ReturnInternal(t *testing.T) {
	ctx := context.Background()

	repo := &fakeUserRepository{
		getUserByEmailErr: errors.New("database down"),
	}

	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	result, err := service.Login(ctx, validLoginRequest("Password123!"))

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.getUserByEmailCalled)
}

func TestLogin_UserNil_ReturnUnauthorized(t *testing.T) {
	ctx := context.Background()

	repo := &fakeUserRepository{
		getUserByEmailResult: nil,
		getUserByEmailErr:    nil,
	}

	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	result, err := service.Login(ctx, validLoginRequest("Password123!"))

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.getUserByEmailCalled)
}

func TestLogin_WrongPassword_ReturnUnauthorized(t *testing.T) {
	ctx := context.Background()

	hash, err := bcrypt.GenerateFromPassword([]byte("CorrectPassword123!"), 12)
	require.NoError(t, err)

	repo := &fakeUserRepository{
		getUserByEmailResult: &authEntities.User{
			Id:       uuid.New(),
			Email:    "giyuu@example.com",
			Password: string(hash),
			FullName: "Giyuu Tomioka",
			CreateAt: time.Now(),
		},
	}

	workspace := &fakeWorkspaceService{}
	service := NewAuthService(repo, workspace)

	result, err := service.Login(ctx, validLoginRequest("WrongPassword123!"))

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.getUserByEmailCalled)
}
