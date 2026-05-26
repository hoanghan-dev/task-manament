package workspace

import (
	"dev/task-management/internal/modules/workspace/dto"
	services "dev/task-management/internal/modules/workspace/services"
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
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		return
	}

	wp, err := h.service.GetWorkspace(c.Request.Context(), ownerId)

	if err != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("Get workspace failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("Get workspace successfully", wp))
}

func (h *WorkspaceHandler) UpdateWorkspace(c *gin.Context) {

	var wpReq dto.WorkspaceUpdateResquestDTO

	err := c.ShouldBindJSON(&wpReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation faild", err.Error()))
		return
	}

	ownerId, err := utils.GetOwnerId(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		return
	}

	updateErr := h.service.UpdateWorkspace(c.Request.Context(), &wpReq, ownerId)

	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", updateErr.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("Update workspace successfully", nil))
}

func (h *WorkspaceHandler) DeleteWorkspace(c *gin.Context) {

	wpIdStr := c.Param("wpId")

	if wpIdStr == "" {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation faild", "workspace id is required"))
		return
	}

	wpId, err := uuid.Parse(wpIdStr)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error",
			"invalid user_id format in context"))
		return
	}

	ownerId, err := utils.GetOwnerId(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ResponseError("Internal Server Error", err.Error()))
		return
	}

	deleteErr := h.service.DeteleWorkspace(c.Request.Context(), wpId, ownerId)

	if deleteErr != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("Workspace not found", deleteErr.Error()))
		return
	}
	c.JSON(http.StatusOK, response.ResponseSuccess("Delete workspace successfully", nil))
}
