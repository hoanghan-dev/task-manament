package handler

import (
	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/services"
	"dev/task-management/pkg/apperror"
	"dev/task-management/pkg/response"
	"dev/task-management/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service services.TaskService
}

func NewTaskHandler(s services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: s,
	}
}

func (h *TaskHandler) GetAllTask(c *gin.Context) {

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	tasks, err := h.service.GetAllTask(c.Request.Context(), ownerId)
	if err != nil {
		apperror.HandleError(c, err) // 404/500 from service/repo
		return
	}

	if tasks == nil {
		c.JSON(http.StatusOK, response.ResponseSuccess("get list task successfully", nil))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("get list task successfully", tasks))
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	idString := c.Param("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		apperror.HandleError(c, apperror.NewValidation("invalid task id format"))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	task, err := h.service.GetTask(c.Request.Context(), id, ownerId)
	if err != nil {
		apperror.HandleError(c, err) // 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("get task successfully", task))
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskReq *request.CreateTaskRequest

	if err := c.ShouldBindJSON(&taskReq); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	task, err := h.service.CreateTask(c.Request.Context(), taskReq, ownerId)
	if err != nil {
		apperror.HandleError(c, err) // 400 Validation / 403 Forbidden / 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusCreated, response.ResponseSuccess("create task successfully", task))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var taskReq *request.UpdateTaskRequest

	if err := c.ShouldBindJSON(&taskReq); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	idString := c.Param("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		apperror.HandleError(c, apperror.NewValidation("invalid task id format"))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	if err := h.service.UpdateTask(c.Request.Context(), id, taskReq, ownerId); err != nil {
		apperror.HandleError(c, err) // 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("update task successfully", nil))
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idString := c.Param("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		apperror.HandleError(c, apperror.NewValidation("invalid task id format"))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	if err := h.service.DeleteTask(c.Request.Context(), id, ownerId); err != nil {
		apperror.HandleError(c, err) // 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("delete task successfully", nil))
}

func (h *TaskHandler) AssignTask(c *gin.Context) {
	var taskReq request.AssignTaskRequest

	if err := c.ShouldBindJSON(&taskReq); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	if err := h.service.AssignTask(c.Request.Context(), &taskReq, ownerId); err != nil {
		apperror.HandleError(c, err) // 403 Forbidden / 404 NotFound / 409 Conflict / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("assign task successfully", nil))
}

func (h *TaskHandler) UpdateTaskStatus(c *gin.Context) {
	var taskReq request.UpdateTaskStatusRequest

	if err := c.ShouldBindJSON(&taskReq); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	taskIdStr := c.Param("id")
	taskId, err := uuid.Parse(taskIdStr)
	if err != nil {
		apperror.HandleError(c, apperror.NewValidation("invalid task id format"))
		return
	}

	ownerId, err := utils.GetOwnerId(c)
	if err != nil {
		apperror.HandleError(c, err) // 401 Unauthorized
		return
	}

	if err := h.service.UpdateTaskStatus(c.Request.Context(), taskId, ownerId, &taskReq); err != nil {
		apperror.HandleError(c, err) // 400 BadRequest / 403 Forbidden / 404 NotFound / 500 Internal
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("update task status successfully", nil))
}
