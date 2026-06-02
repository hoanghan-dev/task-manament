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
	"dev/task-management/pkg/cache"
	"errors"
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

type commentService struct {
	commentRepo repositories.CommentRepository
	redisClient *cache.RedisCacheService
}

func NewCommentService(repo repositories.CommentRepository, redisClient *cache.RedisCacheService) CommentService {
	return &commentService{
		commentRepo: repo,
		redisClient: redisClient,
	}
}

func (s *commentService) CreateComment(ctx context.Context, taskId uuid.UUID, userId uuid.UUID, req *request.CreateCommentRequest) (*response.CommentResponse, error) {
	// Validate content
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, errors.New("comment content cannot be empty")
	}

	// Check task exists
	exists, err := s.commentRepo.TaskExists(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("failed to check task existence: %v", err)
	}
	if !exists {
		return nil, errors.New("task not found")
	}

	// Check user has access to the task
	canAccess, err := s.commentRepo.CanUserAccessTask(ctx, taskId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check user access: %v", err)
	}
	if !canAccess {
		return nil, errors.New("access denied")
	}

	// Create comment entity
	commentId := uuid.New()
	createAt := time.Now()
	comment := entities.NewComment(commentId, taskId, userId, content, createAt)

	// Save to DB
	err = s.commentRepo.Create(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create comment: %v", err)
	}

	// Send WebSocket event asynchronously (don't rollback comment on WS failure)
	go s.notifyCommentCreated(taskId, userId, comment)

	return mapper.EntityToCommentResponse(comment), nil
}

func (s *commentService) GetCommentsByTaskId(ctx context.Context, taskId uuid.UUID, userId uuid.UUID) ([]*response.CommentResponse, error) {
	// Check task exists
	exists, err := s.commentRepo.TaskExists(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("failed to check task existence: %v", err)
	}
	if !exists {
		return nil, errors.New("task not found")
	}

	// Check user has access to the task
	canAccess, err := s.commentRepo.CanUserAccessTask(ctx, taskId, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to check user access: %v", err)
	}
	if !canAccess {
		return nil, errors.New("access denied")
	}

	// Fetch comments
	comments, err := s.commentRepo.FindByTaskId(ctx, taskId)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments: %v", err)
	}

	commentRes := make([]*response.CommentResponse, 0)
	for _, c := range comments {
		commentRes = append(commentRes, mapper.EntityToCommentResponse(&c))
	}

	return commentRes, nil
}

func (s *commentService) DeleteComment(ctx context.Context, commentId uuid.UUID, userId uuid.UUID) error {
	// Check comment exists
	comment, err := s.commentRepo.FindById(ctx, commentId)
	if err != nil {
		return errors.New("comment not found")
	}

	// Check permission: user is comment creator OR workspace owner
	if comment.UserId != userId {
		// Check if user is workspace owner of the task
		canAccess, err := s.commentRepo.CanUserAccessTask(ctx, comment.TaskId, userId)
		if err != nil {
			return fmt.Errorf("failed to check user access: %v", err)
		}
		if !canAccess {
			return errors.New("access denied")
		}
	}

	// Delete from DB
	err = s.commentRepo.Delete(ctx, commentId)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %v", err)
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
