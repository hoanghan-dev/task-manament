package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeTaskRepository struct {
	findAllResult []entities.Task
	findAllErr    error

	findByIDResult *entities.Task
	findByIDErr    error

	createErr           error
	updateErr           error
	deleteErr           error
	assignErr           error
	updateTaskStatusErr error

	taskIsExistsResult          bool
	taskExistsInWorkspaceResult bool

	canUserAccessResult bool
	canUserAccessErr    error

	findAllCalled               int
	findByIDCalled              int
	createCalled                int
	updateCalled                int
	deleteCalled                int
	taskIsExistsCalled          int
	taskExistsInWorkspaceCalled int
	assignCalled                int
	updateTaskStatusCalled      int
	canUserAccessCalled         int

	gotOwnerID     uuid.UUID
	gotTaskID      uuid.UUID
	gotWorkspaceID uuid.UUID
	gotAssigneeID  uuid.UUID
	gotStatus      string
	gotCreatedTask *entities.Task
	gotUpdatedTask *entities.Task
}

func (f *fakeTaskRepository) FindAll(ctx context.Context, ownerId uuid.UUID) ([]entities.Task, error) {
	f.findAllCalled++
	f.gotOwnerID = ownerId
	return f.findAllResult, f.findAllErr
}

func (f *fakeTaskRepository) FindById(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) (*entities.Task, error) {
	f.findByIDCalled++
	f.gotTaskID = id
	f.gotOwnerID = ownerId
	return f.findByIDResult, f.findByIDErr
}

func (f *fakeTaskRepository) CreateTask(ctx context.Context, task *entities.Task) error {
	f.createCalled++
	f.gotCreatedTask = task
	return f.createErr
}

func (f *fakeTaskRepository) UpdateTask(ctx context.Context, id uuid.UUID, task *entities.Task, ownerId uuid.UUID) error {
	f.updateCalled++
	f.gotTaskID = id
	f.gotOwnerID = ownerId
	f.gotUpdatedTask = task
	return f.updateErr
}

func (f *fakeTaskRepository) DeleteTask(ctx context.Context, id uuid.UUID, ownerId uuid.UUID) error {
	f.deleteCalled++
	f.gotTaskID = id
	f.gotOwnerID = ownerId
	return f.deleteErr
}

func (f *fakeTaskRepository) TaskIsExists(ctx context.Context, taskId uuid.UUID) bool {
	f.taskIsExistsCalled++
	f.gotTaskID = taskId
	return f.taskIsExistsResult
}

func (f *fakeTaskRepository) AssignTask(ctx context.Context, taskId uuid.UUID, assigneeId uuid.UUID) error {
	f.assignCalled++
	f.gotTaskID = taskId
	f.gotAssigneeID = assigneeId
	return f.assignErr
}

func (f *fakeTaskRepository) UpdateTaskStatus(ctx context.Context, taskId uuid.UUID, status string) error {
	f.updateTaskStatusCalled++
	f.gotTaskID = taskId
	f.gotStatus = status
	return f.updateTaskStatusErr
}

func (f *fakeTaskRepository) TaskExistsInWorkspace(ctx context.Context, taskId uuid.UUID, workspaceId uuid.UUID) bool {
	f.taskExistsInWorkspaceCalled++
	f.gotTaskID = taskId
	f.gotWorkspaceID = workspaceId
	return f.taskExistsInWorkspaceResult
}

func (f *fakeTaskRepository) CanUserAccessTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (bool, error) {
	f.canUserAccessCalled++
	f.gotTaskID = taskId
	f.gotOwnerID = userId
	return f.canUserAccessResult, f.canUserAccessErr
}

type fakeWorkspaceService struct {
	owned bool
	err   error

	isWorkspaceOwnedByCalled int
	gotWorkspaceID           uuid.UUID
	gotOwnerID               uuid.UUID
}

func (f *fakeWorkspaceService) IsWorkspaceOwnedBy(ctx context.Context, workspaceId uuid.UUID, ownerId uuid.UUID) (bool, error) {
	f.isWorkspaceOwnedByCalled++
	f.gotWorkspaceID = workspaceId
	f.gotOwnerID = ownerId
	return f.owned, f.err
}

type fakeCacheClient struct {
	getErr   error
	setErr   error
	clearErr error
	pushErr  error

	getCalled   int
	setCalled   int
	clearCalled int
	pushCalled  int

	gotGetKey       string
	gotSetKey       string
	gotClearPattern string
	gotPushKey      string
	gotPushValue    any

	pushed chan any
}

func (f *fakeCacheClient) Get(key string, dest any) error {
	f.getCalled++
	f.gotGetKey = key
	return f.getErr
}

func (f *fakeCacheClient) Set(key string, value any, ttl time.Duration) error {
	f.setCalled++
	f.gotSetKey = key
	return f.setErr
}

func (f *fakeCacheClient) Clear(pattern string) error {
	f.clearCalled++
	f.gotClearPattern = pattern
	return f.clearErr
}

func (f *fakeCacheClient) Push(ctx context.Context, key string, value any) error {
	f.pushCalled++
	f.gotPushKey = key
	f.gotPushValue = value
	if f.pushed != nil {
		f.pushed <- value
	}
	return f.pushErr
}

func newTestTaskService(repo *fakeTaskRepository, workspace *fakeWorkspaceService, cache *fakeCacheClient) TaskService {
	if repo == nil {
		repo = &fakeTaskRepository{}
	}
	if workspace == nil {
		workspace = &fakeWorkspaceService{}
	}
	if cache == nil {
		cache = &fakeCacheClient{getErr: errors.New("cache miss")}
	}
	return NewTaskService(repo, workspace, cache)
}

func sampleTask(taskID, workspaceID, assigneeID uuid.UUID, status string) *entities.Task {
	return &entities.Task{
		Id:          taskID,
		Title:       "Test task",
		Description: "Test description",
		Status:      status,
		Assignee:    assigneeID,
		Workspace:   workspaceID,
		CreateAt:    time.Now(),
	}
}

func TestGetAllTask_SuccessFromRepository(t *testing.T) {
	ownerID := uuid.New()
	taskID := uuid.New()
	workspaceID := uuid.New()
	assigneeID := uuid.New()

	repo := &fakeTaskRepository{
		findAllResult: []entities.Task{*sampleTask(taskID, workspaceID, assigneeID, "TODO")},
	}
	cache := &fakeCacheClient{getErr: errors.New("cache miss")}
	service := newTestTaskService(repo, nil, cache)

	result, err := service.GetAllTask(context.Background(), ownerID)

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, taskID, result[0].Id)
	assert.Equal(t, 1, cache.getCalled)
	assert.Equal(t, 1, repo.findAllCalled)
	assert.Equal(t, 1, cache.setCalled)
	assert.Equal(t, "tasks:owner_id:"+ownerID.String(), cache.gotGetKey)
}

func TestGetAllTask_RepositoryError(t *testing.T) {
	repoErr := errors.New("repository error")
	repo := &fakeTaskRepository{findAllErr: repoErr}
	cache := &fakeCacheClient{getErr: errors.New("cache miss")}
	service := newTestTaskService(repo, nil, cache)

	result, err := service.GetAllTask(context.Background(), uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, 1, repo.findAllCalled)
	assert.Equal(t, 0, cache.setCalled)
}

func TestGetAllTask_EmptyList(t *testing.T) {
	repo := &fakeTaskRepository{findAllResult: []entities.Task{}}
	cache := &fakeCacheClient{getErr: errors.New("cache miss")}
	service := newTestTaskService(repo, nil, cache)

	result, err := service.GetAllTask(context.Background(), uuid.New())

	require.NoError(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 1, repo.findAllCalled)
	assert.Equal(t, 0, cache.setCalled)
}

func TestGetAllTask_CacheSetErrorDoesNotFail(t *testing.T) {
	repo := &fakeTaskRepository{findAllResult: []entities.Task{*sampleTask(uuid.New(), uuid.New(), uuid.New(), "TODO")}}
	cache := &fakeCacheClient{getErr: errors.New("cache miss"), setErr: errors.New("cache set failed")}
	service := newTestTaskService(repo, nil, cache)

	result, err := service.GetAllTask(context.Background(), uuid.New())

	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, 1, cache.setCalled)
}

func TestGetTask_SuccessFromRepository(t *testing.T) {
	ownerID := uuid.New()
	taskID := uuid.New()
	workspaceID := uuid.New()
	assigneeID := uuid.New()

	repo := &fakeTaskRepository{
		findByIDResult: sampleTask(taskID, workspaceID, assigneeID, "TODO"),
	}
	cache := &fakeCacheClient{getErr: errors.New("cache miss")}
	service := newTestTaskService(repo, nil, cache)

	result, err := service.GetTask(context.Background(), taskID, ownerID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, taskID, result.Id)
	assert.Equal(t, 1, repo.findByIDCalled)
	assert.Equal(t, 1, cache.setCalled)
}

func TestGetTask_RepositoryError(t *testing.T) {
	repoErr := errors.New("not found")
	repo := &fakeTaskRepository{findByIDErr: repoErr}
	cache := &fakeCacheClient{getErr: errors.New("cache miss")}
	service := newTestTaskService(repo, nil, cache)

	result, err := service.GetTask(context.Background(), uuid.New(), uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, 0, cache.setCalled)
}

func TestCreateTask_Success(t *testing.T) {
	ownerID := uuid.New()
	workspaceID := uuid.New()
	assigneeID := uuid.Nil

	repo := &fakeTaskRepository{}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{getErr: errors.New("cache miss")}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.CreateTaskRequest{
		Title:       "Create test",
		Description: "Description",
		Status:      "TODO",
		Assignee:    assigneeID,
		Workspace:   workspaceID,
	}

	result, err := service.CreateTask(context.Background(), req, ownerID)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, workspace.isWorkspaceOwnedByCalled)
	assert.Equal(t, 1, repo.createCalled)
	assert.Equal(t, workspaceID, repo.gotCreatedTask.Workspace)
	assert.Equal(t, "TODO", repo.gotCreatedTask.Status)
	assert.Equal(t, 1, cache.clearCalled)
	assert.Equal(t, "tasks:*", cache.gotClearPattern)
}

func TestCreateTask_WorkspaceServiceError(t *testing.T) {
	workspaceErr := errors.New("workspace error")
	repo := &fakeTaskRepository{}
	workspace := &fakeWorkspaceService{err: workspaceErr}
	service := newTestTaskService(repo, workspace, nil)

	req := &request.CreateTaskRequest{Workspace: uuid.New(), Status: "TODO"}

	result, err := service.CreateTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, workspaceErr)
	assert.Equal(t, 0, repo.createCalled)
}

func TestCreateTask_AccessDenied(t *testing.T) {
	repo := &fakeTaskRepository{}
	workspace := &fakeWorkspaceService{owned: false}
	service := newTestTaskService(repo, workspace, nil)

	req := &request.CreateTaskRequest{Workspace: uuid.New(), Status: "TODO"}

	result, err := service.CreateTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, repo.createCalled)
}

func TestCreateTask_InvalidStatus(t *testing.T) {
	repo := &fakeTaskRepository{}
	workspace := &fakeWorkspaceService{owned: true}
	service := newTestTaskService(repo, workspace, nil)

	req := &request.CreateTaskRequest{Workspace: uuid.New(), Status: "INVALID_STATUS"}

	result, err := service.CreateTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, 0, repo.createCalled)
}

func TestCreateTask_RepositoryError(t *testing.T) {
	repoErr := errors.New("create failed")
	repo := &fakeTaskRepository{createErr: repoErr}
	workspace := &fakeWorkspaceService{owned: true}
	service := newTestTaskService(repo, workspace, nil)

	req := &request.CreateTaskRequest{Workspace: uuid.New(), Status: "TODO"}

	result, err := service.CreateTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, 1, repo.createCalled)
}

func TestCreateTask_CacheClearErrorDoesNotFail(t *testing.T) {
	repo := &fakeTaskRepository{}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{clearErr: errors.New("clear failed")}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.CreateTaskRequest{Workspace: uuid.New(), Status: "TODO"}

	result, err := service.CreateTask(context.Background(), req, uuid.New())

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestUpdateTask_Success(t *testing.T) {
	taskID := uuid.New()
	ownerID := uuid.New()
	workspaceID := uuid.New()

	repo := &fakeTaskRepository{}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, nil, cache)

	req := &request.UpdateTaskRequest{
		Title:       "Updated title",
		Description: "Updated description",
		Status:      "IN_PROGRESS",
		Workspace:   workspaceID,
	}

	err := service.UpdateTask(context.Background(), taskID, req, ownerID)

	require.NoError(t, err)
	assert.Equal(t, 1, repo.updateCalled)
	assert.Equal(t, taskID, repo.gotTaskID)
	assert.Equal(t, ownerID, repo.gotOwnerID)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestUpdateTask_RepositoryError(t *testing.T) {
	repoErr := errors.New("update failed")
	repo := &fakeTaskRepository{updateErr: repoErr}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, nil, cache)

	err := service.UpdateTask(context.Background(), uuid.New(), &request.UpdateTaskRequest{}, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, 0, cache.clearCalled)
}

func TestUpdateTask_CacheClearErrorDoesNotFail(t *testing.T) {
	repo := &fakeTaskRepository{}
	cache := &fakeCacheClient{clearErr: errors.New("clear failed")}
	service := newTestTaskService(repo, nil, cache)

	err := service.UpdateTask(context.Background(), uuid.New(), &request.UpdateTaskRequest{}, uuid.New())

	require.NoError(t, err)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestDeleteTask_Success(t *testing.T) {
	taskID := uuid.New()
	ownerID := uuid.New()

	repo := &fakeTaskRepository{}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, nil, cache)

	err := service.DeleteTask(context.Background(), taskID, ownerID)

	require.NoError(t, err)
	assert.Equal(t, 1, repo.deleteCalled)
	assert.Equal(t, taskID, repo.gotTaskID)
	assert.Equal(t, ownerID, repo.gotOwnerID)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestDeleteTask_RepositoryError(t *testing.T) {
	repoErr := errors.New("delete failed")
	repo := &fakeTaskRepository{deleteErr: repoErr}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, nil, cache)

	err := service.DeleteTask(context.Background(), uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Equal(t, 0, cache.clearCalled)
}

func TestDeleteTask_CacheClearErrorDoesNotFail(t *testing.T) {
	repo := &fakeTaskRepository{}
	cache := &fakeCacheClient{clearErr: errors.New("clear failed")}
	service := newTestTaskService(repo, nil, cache)

	err := service.DeleteTask(context.Background(), uuid.New(), uuid.New())

	require.NoError(t, err)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestAssignTask_Success(t *testing.T) {
	taskID := uuid.New()
	workspaceID := uuid.New()
	ownerID := uuid.New()
	assigneeID := uuid.New()

	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		taskIsExistsResult:          true,
	}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{pushed: make(chan any, 1)}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{
		TaskId:      taskID,
		WorkspaceId: workspaceID,
		AssigneeId:  assigneeID,
	}

	err := service.AssignTask(context.Background(), req, ownerID)

	require.NoError(t, err)
	assert.Equal(t, 1, repo.taskExistsInWorkspaceCalled)
	assert.Equal(t, 1, repo.taskIsExistsCalled)
	assert.Equal(t, 1, workspace.isWorkspaceOwnedByCalled)
	assert.Equal(t, 1, repo.assignCalled)
	assert.Equal(t, taskID, repo.gotTaskID)
	assert.Equal(t, assigneeID, repo.gotAssigneeID)
	assert.Equal(t, 1, cache.clearCalled)

	select {
	case <-cache.pushed:
		assert.Equal(t, "queue:notifications:assign", cache.gotPushKey)
	case <-time.After(time.Second):
		t.Fatal("expected assign notification event to be pushed")
	}
}

func TestAssignTask_TaskNotInWorkspace(t *testing.T) {
	repo := &fakeTaskRepository{taskExistsInWorkspaceResult: false}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{TaskId: uuid.New(), WorkspaceId: uuid.New(), AssigneeId: uuid.New()}

	err := service.AssignTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.Equal(t, 1, repo.taskExistsInWorkspaceCalled)
	assert.Equal(t, 0, repo.taskIsExistsCalled)
	assert.Equal(t, 0, workspace.isWorkspaceOwnedByCalled)
	assert.Equal(t, 0, repo.assignCalled)
	assert.Equal(t, 0, cache.clearCalled)
}

func TestAssignTask_TaskNotFound(t *testing.T) {
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		taskIsExistsResult:          false,
	}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{TaskId: uuid.New(), WorkspaceId: uuid.New(), AssigneeId: uuid.New()}

	err := service.AssignTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.Equal(t, 1, repo.taskExistsInWorkspaceCalled)
	assert.Equal(t, 1, repo.taskIsExistsCalled)
	assert.Equal(t, 0, workspace.isWorkspaceOwnedByCalled)
	assert.Equal(t, 0, repo.assignCalled)
	assert.Equal(t, 0, cache.clearCalled)
}

func TestAssignTask_WorkspaceServiceError(t *testing.T) {
	workspaceErr := errors.New("workspace error")
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		taskIsExistsResult:          true,
	}
	workspace := &fakeWorkspaceService{err: workspaceErr}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{TaskId: uuid.New(), WorkspaceId: uuid.New(), AssigneeId: uuid.New()}

	err := service.AssignTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, workspaceErr)
	assert.Equal(t, 0, repo.assignCalled)
	assert.Equal(t, 0, cache.clearCalled)
}

func TestAssignTask_AccessDenied(t *testing.T) {
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		taskIsExistsResult:          true,
	}
	workspace := &fakeWorkspaceService{owned: false}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{TaskId: uuid.New(), WorkspaceId: uuid.New(), AssigneeId: uuid.New()}

	err := service.AssignTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.Equal(t, 1, workspace.isWorkspaceOwnedByCalled)
	assert.Equal(t, 0, repo.assignCalled)
	assert.Equal(t, 0, cache.clearCalled)
}

func TestAssignTask_RepositoryAssignError(t *testing.T) {
	assignErr := errors.New("assign failed")
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		taskIsExistsResult:          true,
		assignErr:                   assignErr,
	}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{TaskId: uuid.New(), WorkspaceId: uuid.New(), AssigneeId: uuid.New()}

	err := service.AssignTask(context.Background(), req, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, assignErr)
	assert.Equal(t, 1, repo.assignCalled)
	assert.Equal(t, 0, cache.clearCalled)
	assert.Equal(t, 0, cache.pushCalled)
}

func TestAssignTask_CacheClearErrorDoesNotFail(t *testing.T) {
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		taskIsExistsResult:          true,
	}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{clearErr: errors.New("clear failed"), pushed: make(chan any, 1)}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{TaskId: uuid.New(), WorkspaceId: uuid.New(), AssigneeId: uuid.New()}

	err := service.AssignTask(context.Background(), req, uuid.New())

	require.NoError(t, err)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestAssignTask_NotificationPushErrorDoesNotFail(t *testing.T) {
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		taskIsExistsResult:          true,
	}
	workspace := &fakeWorkspaceService{owned: true}
	cache := &fakeCacheClient{pushErr: errors.New("push failed"), pushed: make(chan any, 1)}
	service := newTestTaskService(repo, workspace, cache)

	req := &request.AssignTaskRequest{TaskId: uuid.New(), WorkspaceId: uuid.New(), AssigneeId: uuid.New()}

	err := service.AssignTask(context.Background(), req, uuid.New())

	require.NoError(t, err)

	select {
	case <-cache.pushed:
		assert.Equal(t, "queue:notifications:assign", cache.gotPushKey)
	case <-time.After(time.Second):
		t.Fatal("expected assign notification push to be attempted")
	}
}

func TestUpdateTaskStatus_Success(t *testing.T) {
	taskID := uuid.New()
	ownerID := uuid.New()
	workspaceID := uuid.New()
	assigneeID := uuid.New()

	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		findByIDResult:              sampleTask(taskID, workspaceID, assigneeID, "TODO"),
		canUserAccessResult:         true,
	}
	cache := &fakeCacheClient{pushed: make(chan any, 1)}
	service := newTestTaskService(repo, nil, cache)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: workspaceID, Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), taskID, ownerID, req)

	require.NoError(t, err)
	assert.Equal(t, 1, repo.taskExistsInWorkspaceCalled)
	assert.Equal(t, 1, repo.canUserAccessCalled)
	assert.Equal(t, 1, repo.updateTaskStatusCalled)
	assert.Equal(t, "IN_PROGRESS", repo.gotStatus)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestUpdateTaskStatus_TaskNotInWorkspace(t *testing.T) {
	repo := &fakeTaskRepository{taskExistsInWorkspaceResult: false}
	cache := &fakeCacheClient{}
	service := newTestTaskService(repo, nil, cache)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: uuid.New(), Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), uuid.New(), uuid.New(), req)

	require.Error(t, err)
	assert.Equal(t, 1, repo.taskExistsInWorkspaceCalled)
	assert.Equal(t, 0, repo.findByIDCalled)
	assert.Equal(t, 0, repo.updateTaskStatusCalled)
}

func TestUpdateTaskStatus_FindByIDError(t *testing.T) {
	findErr := errors.New("find failed")
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		findByIDErr:                 findErr,
	}
	service := newTestTaskService(repo, nil, nil)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: uuid.New(), Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), uuid.New(), uuid.New(), req)

	require.Error(t, err)
	assert.ErrorIs(t, err, findErr)
	assert.Equal(t, 0, repo.canUserAccessCalled)
	assert.Equal(t, 0, repo.updateTaskStatusCalled)
}

func TestUpdateTaskStatus_CanUserAccessError(t *testing.T) {
	accessErr := errors.New("access check failed")
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		findByIDResult:              sampleTask(uuid.New(), uuid.New(), uuid.New(), "TODO"),
		canUserAccessErr:            accessErr,
	}
	service := newTestTaskService(repo, nil, nil)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: uuid.New(), Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), uuid.New(), uuid.New(), req)

	require.Error(t, err)
	assert.ErrorIs(t, err, accessErr)
	assert.Equal(t, 0, repo.updateTaskStatusCalled)
}

func TestUpdateTaskStatus_AccessDenied(t *testing.T) {
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		findByIDResult:              sampleTask(uuid.New(), uuid.New(), uuid.New(), "TODO"),
		canUserAccessResult:         false,
	}
	service := newTestTaskService(repo, nil, nil)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: uuid.New(), Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), uuid.New(), uuid.New(), req)

	require.Error(t, err)
	assert.Equal(t, 0, repo.updateTaskStatusCalled)
}

func TestUpdateTaskStatus_InvalidTransition(t *testing.T) {
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		findByIDResult:              sampleTask(uuid.New(), uuid.New(), uuid.New(), "DONE"),
		canUserAccessResult:         true,
	}
	service := newTestTaskService(repo, nil, nil)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: uuid.New(), Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), uuid.New(), uuid.New(), req)

	require.Error(t, err)
	assert.Equal(t, 0, repo.updateTaskStatusCalled)
}

func TestUpdateTaskStatus_RepositoryUpdateError(t *testing.T) {
	updateErr := errors.New("update status failed")
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		findByIDResult:              sampleTask(uuid.New(), uuid.New(), uuid.New(), "TODO"),
		canUserAccessResult:         true,
		updateTaskStatusErr:         updateErr,
	}
	service := newTestTaskService(repo, nil, nil)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: uuid.New(), Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), uuid.New(), uuid.New(), req)

	require.Error(t, err)
	assert.ErrorIs(t, err, updateErr)
}

func TestUpdateTaskStatus_CacheClearErrorDoesNotFail(t *testing.T) {
	repo := &fakeTaskRepository{
		taskExistsInWorkspaceResult: true,
		findByIDResult:              sampleTask(uuid.New(), uuid.New(), uuid.New(), "TODO"),
		canUserAccessResult:         true,
	}
	cache := &fakeCacheClient{clearErr: errors.New("clear failed"), pushed: make(chan any, 1)}
	service := newTestTaskService(repo, nil, cache)

	req := &request.UpdateTaskStatusRequest{WorkspaceId: uuid.New(), Status: "IN_PROGRESS"}

	err := service.UpdateTaskStatus(context.Background(), uuid.New(), uuid.New(), req)

	require.NoError(t, err)
	assert.Equal(t, 1, cache.clearCalled)
}

func TestChangeStatusIsValid(t *testing.T) {
	service := newTestTaskService(nil, nil, nil).(*taskService)

	tests := []struct {
		name      string
		oldStatus string
		newStatus string
		expected  bool
	}{
		{name: "TODO to IN_PROGRESS", oldStatus: "TODO", newStatus: "IN_PROGRESS", expected: true},
		{name: "IN_PROGRESS to DONE", oldStatus: "IN_PROGRESS", newStatus: "DONE", expected: true},
		{name: "same status", oldStatus: "TODO", newStatus: "TODO", expected: false},
		{name: "DONE cannot change", oldStatus: "DONE", newStatus: "IN_PROGRESS", expected: false},
		{name: "cannot change back to TODO", oldStatus: "IN_PROGRESS", newStatus: "TODO", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := service.changeStatusIsValid(tt.oldStatus, tt.newStatus)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
