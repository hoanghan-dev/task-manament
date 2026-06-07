package services

import (
	"context"
	"dev/task-management/internal/modules/comment/dto/request"
	"dev/task-management/internal/modules/comment/dto/response"
	"dev/task-management/internal/modules/comment/entities"
	"dev/task-management/internal/modules/comment/events"
	"dev/task-management/internal/modules/comment/mapper"
	"dev/task-management/internal/modules/comment/repositories"
	"dev/task-management/internal/realtime"
	"dev/task-management/pkg/apperror"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CommentService interface {
	CreateComment(ctx context.Context, taskId uuid.UUID, userId uuid.UUID, req *request.CreateCommentRequest) (*response.CommentResponse, error)
	GetCommentsByTaskId(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) ([]*response.CommentResponse, error)
	DeleteComment(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) error
}

type CommentPublisher interface {
	Publish(ctx context.Context, channel string, value any) error
}

type commentService struct {
	commentRepo repositories.CommentRepository
	redisClient CommentPublisher
}

func NewCommentService(repo repositories.CommentRepository, redisClient CommentPublisher) CommentService {
	return &commentService{
		commentRepo: repo,
		redisClient: redisClient,
	}
}

func (s *commentService) CreateComment(ctx context.Context, taskId uuid.UUID, userId uuid.UUID, req *request.CreateCommentRequest) (*response.CommentResponse, error) {
	// Validate content
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, apperror.NewValidation("comment content cannot be empty")
	}

	// Check task exists
	exists, err := s.commentRepo.TaskExists(ctx, taskId)
	if err != nil {
		return nil, err // AppError from repository
	}
	if !exists {
		return nil, apperror.NewNotFound("task")
	}

	// Check user has access to the task
	canAccess, err := s.commentRepo.CanUserAccessTask(ctx, taskId, userId)
	if err != nil {
		return nil, err // AppError from repository
	}
	if !canAccess {
		return nil, apperror.NewForbidden("you don't have permission to comment on this task")
	}

	// Create comment entity
	commentId := uuid.New()
	createAt := time.Now()
	comment := entities.NewComment(commentId, taskId, userId, content, createAt)

	// Save to DB
	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, err // AppError from repository
	}

	// Send WebSocket event asynchronously (don't rollback comment on WS failure)
	go s.notifyCommentCreated(taskId, userId, comment)

	return mapper.EntityToCommentResponse(comment), nil
}

func (s *commentService) GetCommentsByTaskId(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) ([]*response.CommentResponse, error) {
	// Check task exists
	exists, err := s.commentRepo.TaskExists(ctx, taskId)
	if err != nil {
		return nil, err // AppError from repository
	}
	if !exists {
		return nil, apperror.NewNotFound("task")
	}

	// Check user has access to the task
	canAccess, err := s.commentRepo.CanUserAccessTask(ctx, taskId, userId)
	if err != nil {
		return nil, err // AppError from repository
	}
	if !canAccess {
		return nil, apperror.NewForbidden("you don't have permission to view comments on this task")
	}

	// Fetch comments
	comments, err := s.commentRepo.FindByTaskId(ctx, taskId)
	if err != nil {
		return nil, err // AppError from repository
	}

	commentRes := make([]*response.CommentResponse, 0)
	for _, c := range comments {
		commentRes = append(commentRes, mapper.EntityToCommentResponse(&c))
	}

	return commentRes, nil
}

func (s *commentService) DeleteComment(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) error {
	// Check comment exists - repository returns AppError(NotFound) or AppError(Internal)
	comment, err := s.commentRepo.FindById(ctx, commentId)
	if err != nil {
		return err // AppError: 404 if not found, 500 if DB error
	}

	// Check permission: user is comment creator OR workspace owner
	if comment.UserId != userId {
		canAccess, err := s.commentRepo.CanUserAccessTask(ctx, comment.TaskId, userId)
		if err != nil {
			return err // AppError from repository
		}
		if !canAccess {
			return apperror.NewForbidden("you don't have permission to delete this comment")
		}
	}

	// Delete from DB
	err = s.commentRepo.Delete(ctx, commentId)
	if err != nil {
		return err // AppError from repository
	}

	// Send WebSocket event asynchronously
	go s.notifyCommentDeleted(comment.TaskId, userId, comment)

	return nil
}

func (s *commentService) notifyCommentCreated(taskId uuid.UUID, senderId uuid.UUID, comment *entities.Comment) {
	ctx := context.Background()

	receivers, err := s.commentRepo.FindCommentReceiversByTaskId(ctx, taskId)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get comment receivers: %v\n", err)
		return
	}

	event := &realtime.Event{
		Type: "comment.created",
		Data: events.NewCommentCreatedEvent(comment.CommentId, comment.TaskId, comment.UserId, comment.Content, comment.CreateAt),
	}

	for _, receiverId := range receivers {
		// Skip sender to avoid echo
		if receiverId == senderId {
			continue
		}

		channel := fmt.Sprintf("notify:comment:%s", receiverId.String())
		err := s.redisClient.Publish(ctx, channel, event)
		if err != nil {
			fmt.Printf("[ERROR] Failed to publish comment.created event to user %s: %v\n", receiverId.String(), err)
		}
	}
}

func (s *commentService) notifyCommentDeleted(taskId uuid.UUID, senderId uuid.UUID, comment *entities.Comment) {
	ctx := context.Background()

	receivers, err := s.commentRepo.FindCommentReceiversByTaskId(ctx, taskId)
	if err != nil {
		fmt.Printf("[ERROR] Failed to get comment receivers: %v\n", err)
		return
	}

	event := &realtime.Event{
		Type: "comment.deleted",
		Data: events.NewCommentDeletedEvent(comment.CommentId, comment.TaskId, comment.UserId),
	}

	for _, receiverId := range receivers {
		if receiverId == senderId {
			continue
		}

		channel := fmt.Sprintf("notify:comment:%s", receiverId.String())
		err := s.redisClient.Publish(ctx, channel, event)
		if err != nil {
			fmt.Printf("[ERROR] Failed to publish comment.deleted event to user %s: %v\n", receiverId.String(), err)
		}
	}
}
