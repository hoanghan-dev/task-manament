package mapper

import (
	"dev/task-management/internal/modules/comment/dto/response"
	"dev/task-management/internal/modules/comment/entities"
)

func EntityToCommentResponse(c *entities.Comment) *response.CommentResponse {
	return &response.CommentResponse{
		CommentId: c.CommentId,
		TaskId:    c.TaskId,
		UserId:    c.UserId,
		Content:   c.Content,
		CreateAt:  c.CreateAt,
		UpdateAt:  c.UpdateAt,
	}
}
