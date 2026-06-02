package entities

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	CommentId uuid.UUID
	TaskId    uuid.UUID
	UserId    uuid.UUID
	Content   string
	CreateAt  time.Time
	UpdateAt  *time.Time
}

func NewComment(commentId uuid.UUID, taskId uuid.UUID, userId uuid.UUID, content string, createAt time.Time) *Comment {
	return &Comment{
		CommentId: commentId,
		TaskId:    taskId,
		UserId:    userId,
		Content:   content,
		CreateAt:  createAt,
		UpdateAt:  nil,
	}
}
