package workspace

import (
	"dev/task-management/internal/modules/workspace/dto"
	services "dev/task-management/internal/modules/workspace/services"
	"dev/task-management/pkg/apperror"
	"dev/task-management/pkg/response"
	"dev/task-management/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkspaceHandler struct {
	service services.WorkspaceService
}

func NewWorkspaceHandler(service services.WorkspaceService) *WorkspaceHandler {
	return &WorkspaceHandler{
		service: service,
	}
}

func (h *WorkspaceHandler) GetWorkspace(c *gin.Context) {
	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	wp, err := h.service.GetWorkspace(c.Request.Context(), ownerId)
	if err != nil {
		apperror.HandleError(c, err) // 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("get workspace successfully", wp))
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {
	var wpReq dto.WorkspaceUpdateResquestDTO

	if err := c.ShouldBindJSON(&wpReq); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	if err := h.service.UpdateWorkspace(c.Request.Context(), &wpReq, ownerId); err != nil {
		apperror.HandleError(c, err) // 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("update workspace successfully", nil))
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {
	wpIdStr := c.Param("wpId")

	if wpIdStr == "" {
		apperror.HandleError(c, apperror.NewValidation("workspace id is required"))
		return
	}

	wpId, err := uuid.Parse(wpIdStr)
	if err != nil {
		// was incorrectly returning 500 — this is clearly a 400 validation error
		apperror.HandleError(c, apperror.NewValidation("invalid workspace id format"))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	if err := h.service.DeteleWorkspace(c.Request.Context(), wpId, ownerId); err != nil {
		apperror.HandleError(c, err) // 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("delete workspace successfully", nil))
}
