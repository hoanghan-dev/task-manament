package handler

import (
	"dev/task-management/internal/modules/comment/dto/request"
	"dev/task-management/internal/modules/comment/services"
	"dev/task-management/pkg/apperror"
	"dev/task-management/pkg/response"
	"dev/task-management/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CommentHandler struct {
	service services.CommentService
}

func NewCommentHandler(s services.CommentService) *CommentHandler {
	return &CommentHandler{
		service: s,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	taskIdStr := c.Param("id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		apperror.HandleError(c, apperror.NewValidation("invalid task id format"))
		return
	}

	var req request.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	userId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	comment, err := h.service.CreateComment(c.Request.Context(), taskId, userId, &req)
	if err != nil {
		// 400 Validation / 403 Forbidden / 404 NotFound / 500 Internal
		// No more fragile switch err.Error() — AppError drives the mapping
		apperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.ResponseSuccess("create comment successfully", comment))
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	taskIdStr := c.Param("id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		apperror.HandleError(c, apperror.NewValidation("invalid task id format"))
		return
	}

	userId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	comments, err := h.service.GetCommentsByTaskId(c.Request.Context(), taskId, userId)
	if err != nil {
		// 403 Forbidden / 404 NotFound / 500 Internal
		apperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("get comments successfully", comments))
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIdStr := c.Param("commentId")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		apperror.HandleError(c, apperror.NewValidation("invalid comment id format"))
		return
	}

	userId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	if err := h.service.DeleteComment(c.Request.Context(), commentId, userId); err != nil {
		// 403 Forbidden / 404 NotFound / 500 Internal
		apperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("delete comment successfully", nil))
}
