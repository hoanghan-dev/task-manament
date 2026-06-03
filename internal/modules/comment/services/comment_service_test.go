package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"dev/task-management/internal/modules/comment/dto/request"
	"dev/task-management/internal/modules/comment/entities"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeCommentRepository struct {
	createCalled bool
	createArg    *entities.Comment
	createErr    error

	findByTaskIdCalled bool
	findByTaskIdArg    uuid.UUID
	findByTaskIdResult []entities.Comment
	findByTaskIdErr    error

	findByIdCalled bool
	findByIdArg    uuid.UUID
	findByIdResult *entities.Comment
	findByIdErr    error

	deleteCalled bool
	deleteArg    uuid.UUID
	deleteErr    error

	taskExistsCalled bool
	taskExistsArg    uuid.UUID
	taskExistsResult bool
	taskExistsErr    error

	canUserAccessTaskCalled bool
	canUserAccessTaskTaskId uuid.UUID
	canUserAccessTaskUserId uuid.UUID
	canUserAccessTaskResult bool
	canUserAccessTaskErr    error

	findReceiversCalled bool
	findReceiversTaskId uuid.UUID
	findReceiversResult []uuid.UUID
	findReceiversErr    error
}

func (f *fakeCommentRepository) Create(ctx context.Context, comment *entities.Comment) error {
	f.createCalled = true
	f.createArg = comment
	return f.createErr
}

func (f *fakeCommentRepository) FindByTaskId(ctx context.Context, taskId uuid.UUID) ([]entities.Comment, error) {
	f.findByTaskIdCalled = true
	f.findByTaskIdArg = taskId
	return f.findByTaskIdResult, f.findByTaskIdErr
}

func (f *fakeCommentRepository) FindById(ctx context.Context, commentId uuid.UUID) (*entities.Comment, error) {
	f.findByIdCalled = true
	f.findByIdArg = commentId
	return f.findByIdResult, f.findByIdErr
}

func (f *fakeCommentRepository) Delete(ctx context.Context, commentId uuid.UUID) error {
	f.deleteCalled = true
	f.deleteArg = commentId
	return f.deleteErr
}

func (f *fakeCommentRepository) TaskExists(ctx context.Context, taskId uuid.UUID) (bool, error) {
	f.taskExistsCalled = true
	f.taskExistsArg = taskId
	return f.taskExistsResult, f.taskExistsErr
}

func (f *fakeCommentRepository) CanUserAccessTask(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) (bool, error) {
	f.canUserAccessTaskCalled = true
	f.canUserAccessTaskTaskId = taskId
	f.canUserAccessTaskUserId = userId
	return f.canUserAccessTaskResult, f.canUserAccessTaskErr
}

func (f *fakeCommentRepository) FindCommentReceiversByTaskId(ctx context.Context, taskId uuid.UUID) ([]uuid.UUID, error) {
	f.findReceiversCalled = true
	f.findReceiversTaskId = taskId
	return f.findReceiversResult, f.findReceiversErr
}

type fakeCommentPublisher struct {
	publishCalled bool
	publishCount  int
	channels      []string
	publishErr    error
}

func (f *fakeCommentPublisher) Publish(ctx context.Context, channel string, value any) error {
	f.publishCalled = true
	f.publishCount++
	f.channels = append(f.channels, channel)
	return f.publishErr
}

func TestCreateComment_Success(t *testing.T) {
	ctx := context.Background()
	taskId := uuid.New()
	userId := uuid.New()

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskResult: true,
	}
	publisher := &fakeCommentPublisher{}

	service := NewCommentService(repo, publisher)

	result, err := service.CreateComment(ctx, taskId, userId, &request.CreateCommentRequest{
		Content: "  hello comment  ",
	})

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.True(t, repo.taskExistsCalled)
	assert.Equal(t, taskId, repo.taskExistsArg)

	assert.True(t, repo.canUserAccessTaskCalled)
	assert.Equal(t, taskId, repo.canUserAccessTaskTaskId)
	assert.Equal(t, userId, repo.canUserAccessTaskUserId)

	assert.True(t, repo.createCalled)
	require.NotNil(t, repo.createArg)
	assert.Equal(t, taskId, repo.createArg.TaskId)
	assert.Equal(t, userId, repo.createArg.UserId)
	assert.Equal(t, "hello comment", repo.createArg.Content)
}

func TestCreateComment_EmptyContent(t *testing.T) {
	ctx := context.Background()

	repo := &fakeCommentRepository{}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.CreateComment(ctx, uuid.New(), uuid.New(), &request.CreateCommentRequest{
		Content: "   ",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.False(t, repo.taskExistsCalled)
	assert.False(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.createCalled)
}

func TestCreateComment_TaskExistsRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("task exists db error")

	repo := &fakeCommentRepository{
		taskExistsErr: expectedErr,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.CreateComment(ctx, uuid.New(), uuid.New(), &request.CreateCommentRequest{
		Content: "hello",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.True(t, repo.taskExistsCalled)
	assert.False(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.createCalled)
}

func TestCreateComment_TaskNotFound(t *testing.T) {
	ctx := context.Background()

	repo := &fakeCommentRepository{
		taskExistsResult: false,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.CreateComment(ctx, uuid.New(), uuid.New(), &request.CreateCommentRequest{
		Content: "hello",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.taskExistsCalled)
	assert.False(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.createCalled)
}

func TestCreateComment_CanAccessRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("access db error")

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskErr:    expectedErr,
		canUserAccessTaskResult: false,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.CreateComment(ctx, uuid.New(), uuid.New(), &request.CreateCommentRequest{
		Content: "hello",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.True(t, repo.taskExistsCalled)
	assert.True(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.createCalled)
}

func TestCreateComment_AccessDenied(t *testing.T) {
	ctx := context.Background()

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskResult: false,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.CreateComment(ctx, uuid.New(), uuid.New(), &request.CreateCommentRequest{
		Content: "hello",
	})

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.taskExistsCalled)
	assert.True(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.createCalled)
}

func TestCreateComment_CreateRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("create comment db error")

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskResult: true,
		createErr:               expectedErr,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.CreateComment(ctx, uuid.New(), uuid.New(), &request.CreateCommentRequest{
		Content: "hello",
	})

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.True(t, repo.createCalled)
}

func TestGetCommentsByTaskId_Success(t *testing.T) {
	ctx := context.Background()
	taskId := uuid.New()
	userId := uuid.New()
	commentId := uuid.New()

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskResult: true,
		findByTaskIdResult: []entities.Comment{
			{
				CommentId: commentId,
				TaskId:    taskId,
				UserId:    userId,
				Content:   "hello",
				CreateAt:  time.Now(),
			},
		},
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.GetCommentsByTaskId(ctx, taskId, userId)

	require.NoError(t, err)
	require.Len(t, result, 1)

	assert.True(t, repo.taskExistsCalled)
	assert.True(t, repo.canUserAccessTaskCalled)
	assert.True(t, repo.findByTaskIdCalled)
	assert.Equal(t, taskId, repo.findByTaskIdArg)
}

func TestGetCommentsByTaskId_TaskExistsRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("task exists db error")

	repo := &fakeCommentRepository{
		taskExistsErr: expectedErr,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.GetCommentsByTaskId(ctx, uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.True(t, repo.taskExistsCalled)
	assert.False(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.findByTaskIdCalled)
}

func TestGetCommentsByTaskId_TaskNotFound(t *testing.T) {
	ctx := context.Background()

	repo := &fakeCommentRepository{
		taskExistsResult: false,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.GetCommentsByTaskId(ctx, uuid.New(), uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.taskExistsCalled)
	assert.False(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.findByTaskIdCalled)
}

func TestGetCommentsByTaskId_AccessDenied(t *testing.T) {
	ctx := context.Background()

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskResult: false,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.GetCommentsByTaskId(ctx, uuid.New(), uuid.New())

	require.Error(t, err)
	assert.Nil(t, result)
	assert.True(t, repo.taskExistsCalled)
	assert.True(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.findByTaskIdCalled)
}

func TestGetCommentsByTaskId_FindRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("find comments db error")

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskResult: true,
		findByTaskIdErr:         expectedErr,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.GetCommentsByTaskId(ctx, uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.True(t, repo.findByTaskIdCalled)
}

func TestGetCommentsByTaskId_EmptyList(t *testing.T) {
	ctx := context.Background()

	repo := &fakeCommentRepository{
		taskExistsResult:        true,
		canUserAccessTaskResult: true,
		findByTaskIdResult:      []entities.Comment{},
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	result, err := service.GetCommentsByTaskId(ctx, uuid.New(), uuid.New())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func TestDeleteComment_Success_ByCommentCreator(t *testing.T) {
	ctx := context.Background()
	commentId := uuid.New()
	taskId := uuid.New()
	userId := uuid.New()

	repo := &fakeCommentRepository{
		findByIdResult: &entities.Comment{
			CommentId: commentId,
			TaskId:    taskId,
			UserId:    userId,
			Content:   "hello",
			CreateAt:  time.Now(),
		},
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	err := service.DeleteComment(ctx, commentId, userId)

	require.NoError(t, err)
	assert.True(t, repo.findByIdCalled)
	assert.Equal(t, commentId, repo.findByIdArg)
	assert.False(t, repo.canUserAccessTaskCalled, "creator deletes own comment, no access check needed")
	assert.True(t, repo.deleteCalled)
	assert.Equal(t, commentId, repo.deleteArg)
}

func TestDeleteComment_Success_ByTaskAccessibleUser(t *testing.T) {
	ctx := context.Background()
	commentId := uuid.New()
	taskId := uuid.New()
	commentCreatorId := uuid.New()
	ownerOrAssigneeId := uuid.New()

	repo := &fakeCommentRepository{
		findByIdResult: &entities.Comment{
			CommentId: commentId,
			TaskId:    taskId,
			UserId:    commentCreatorId,
			Content:   "hello",
			CreateAt:  time.Now(),
		},
		canUserAccessTaskResult: true,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	err := service.DeleteComment(ctx, commentId, ownerOrAssigneeId)

	require.NoError(t, err)
	assert.True(t, repo.findByIdCalled)
	assert.True(t, repo.canUserAccessTaskCalled)
	assert.Equal(t, taskId, repo.canUserAccessTaskTaskId)
	assert.Equal(t, ownerOrAssigneeId, repo.canUserAccessTaskUserId)
	assert.True(t, repo.deleteCalled)
}

func TestDeleteComment_FindByIdRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("find comment db error")

	repo := &fakeCommentRepository{
		findByIdErr: expectedErr,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	err := service.DeleteComment(ctx, uuid.New(), uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, repo.findByIdCalled)
	assert.False(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.deleteCalled)
}

func TestDeleteComment_CanAccessRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("access db error")

	commentId := uuid.New()
	taskId := uuid.New()

	repo := &fakeCommentRepository{
		findByIdResult: &entities.Comment{
			CommentId: commentId,
			TaskId:    taskId,
			UserId:    uuid.New(),
			Content:   "hello",
			CreateAt:  time.Now(),
		},
		canUserAccessTaskErr: expectedErr,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	err := service.DeleteComment(ctx, commentId, uuid.New())

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, repo.findByIdCalled)
	assert.True(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.deleteCalled)
}

func TestDeleteComment_AccessDenied(t *testing.T) {
	ctx := context.Background()

	commentId := uuid.New()
	taskId := uuid.New()

	repo := &fakeCommentRepository{
		findByIdResult: &entities.Comment{
			CommentId: commentId,
			TaskId:    taskId,
			UserId:    uuid.New(),
			Content:   "hello",
			CreateAt:  time.Now(),
		},
		canUserAccessTaskResult: false,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	err := service.DeleteComment(ctx, commentId, uuid.New())

	require.Error(t, err)
	assert.True(t, repo.findByIdCalled)
	assert.True(t, repo.canUserAccessTaskCalled)
	assert.False(t, repo.deleteCalled)
}

func TestDeleteComment_DeleteRepositoryError(t *testing.T) {
	ctx := context.Background()
	expectedErr := errors.New("delete db error")

	commentId := uuid.New()
	userId := uuid.New()

	repo := &fakeCommentRepository{
		findByIdResult: &entities.Comment{
			CommentId: commentId,
			TaskId:    uuid.New(),
			UserId:    userId,
			Content:   "hello",
			CreateAt:  time.Now(),
		},
		deleteErr: expectedErr,
	}
	publisher := &fakeCommentPublisher{}
	service := NewCommentService(repo, publisher)

	err := service.DeleteComment(ctx, commentId, userId)

	require.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
	assert.True(t, repo.deleteCalled)
}
