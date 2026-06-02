package handler

import (
	"dev/task-management/internal/modules/comment/dto/request"
	"dev/task-management/internal/modules/comment/services"
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
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", "invalid task id format"))
		return
	}

	var req request.CreateCommentRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		return
	}

	userId, err := utils.GetOwnerId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		return
	}

	comment, err := h.service.CreateComment(c.Request.Context(), taskId, userId, &req)
	if err != nil {
		switch err.Error() {
		case "task not found":
			c.JSON(http.StatusNotFound, response.ResponseError("create comment failed", err.Error()))
		case "access denied":
			c.JSON(http.StatusForbidden, response.ResponseError("create comment failed", err.Error()))
		case "comment content cannot be empty":
			c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		}
		return
	}

	c.JSON(http.StatusCreated, response.ResponseSuccess("create comment successfully", comment))
}

func (h *CommentHandler) GetComments(c *gin.Context) {
	taskIdStr := c.Param("id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", "invalid task id format"))
		return
	}

	userId, err := utils.GetOwnerId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		return
	}

	comments, err := h.service.GetCommentsByTaskId(c.Request.Context(), taskId, userId)
	if err != nil {
		switch err.Error() {
		case "task not found":
			c.JSON(http.StatusNotFound, response.ResponseError("get comments failed", err.Error()))
		case "access denied":
			c.JSON(http.StatusForbidden, response.ResponseError("get comments failed", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("get comments successfully", comments))
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentIdStr := c.Param("commentId")
	commentId, err := uuid.Parse(commentIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", "invalid comment id format"))
		return
	}

	userId, err := utils.GetOwnerId(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		return
	}

	err = h.service.DeleteComment(c.Request.Context(), commentId, userId)
	if err != nil {
		switch err.Error() {
		case "comment not found":
			c.JSON(http.StatusNotFound, response.ResponseError("delete comment failed", err.Error()))
		case "access denied":
			c.JSON(http.StatusForbidden, response.ResponseError("delete comment failed", err.Error()))
		default:
			c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		}
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("delete comment successfully", nil))
}
