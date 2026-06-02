package response

import (
	"time"

	"github.com/google/uuid"
)

type CommentResponse struct {
	CommentId uuid.UUID  `json:"comment_id"`
	TaskId    uuid.UUID  `json:"task_id"`
	UserId    uuid.UUID  `json:"user_id"`
	Content   string     `json:"content"`
	CreateAt  time.Time  `json:"create_at"`
	UpdateAt  *time.Time `json:"update_at,omitempty"`
}
