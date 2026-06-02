package events

import (
	"time"

	"github.com/google/uuid"
)

type CommentCreatedEvent struct {
	CommentId uuid.UUID `json:"comment_id"`
	TaskId    uuid.UUID `json:"task_id"`
	UserId    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreateAt  time.Time `json:"create_at"`
}

func NewCommentCreatedEvent(commentId uuid.UUID, taskId uuid.UUID, userId uuid.UUID, content string, createAt time.Time) *CommentCreatedEvent {
	return &CommentCreatedEvent{
		CommentId: commentId,
		TaskId:    taskId,
		UserId:    userId,
		Content:   content,
		CreateAt:  createAt,
	}
}

type CommentDeletedEvent struct {
	CommentId uuid.UUID `json:"comment_id"`
	TaskId    uuid.UUID `json:"task_id"`
	UserId    uuid.UUID `json:"user_id"`
}

func NewCommentDeletedEvent(commentId uuid.UUID, taskId uuid.UUID, userId uuid.UUID) *CommentDeletedEvent {
	return &CommentDeletedEvent{
		CommentId: commentId,
		TaskId:    taskId,
		UserId:    userId,
	}
}
