package handler

import (
	"dev/task-management/internal/modules/task/dto/request"
	"dev/task-management/internal/modules/task/services"
	"dev/task-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(s *services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: s,
	}
}

func (h *TaskHandler) GetAllTask(c *gin.Context) {
	tasks, err := h.service.GetAllTask(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("get list task failed", err.Error()))
		return
	}
	c.JSON(200, response.ResponseSuccess("get list task successfully", tasks))
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	idString := c.Param("id")
	id, err := uuid.Parse(idString)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		return
	}
	task, err1 := h.service.GetTask(c.Request.Context(), id)

	if err1 != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("get task failed", err1.Error()))
		return
	}

	c.JSON(200, response.ResponseSuccess("get task successfully", task))
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskReq *request.TaskRequest

	err := c.ShouldBindJSON(&taskReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		return
	}
	task, err := h.service.CreateTask(c.Request.Context(), taskReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("error create task", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.ResponseSuccess("create task successfully ", task))
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var taskReq *request.TaskRequest

	err := c.ShouldBindJSON(&taskReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		return
	}

	idString := c.Param("id")

	id, err := uuid.Parse(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		return
	}
	err2 := h.service.UpdateTask(c.Request.Context(), id, taskReq)

	if err2 != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err2.Error()))
		return
	}
	c.JSON(http.StatusOK, response.ResponseSuccess("update task successfully", nil))
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idString := c.Param("id")

	id, err := uuid.Parse(idString)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		return
	}
	err2 := h.service.DeleteTask(c.Request.Context(), id)
	if err2 != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("delete task failed", err2.Error()))
		return
	}
	c.JSON(http.StatusOK, response.ResponseSuccess("delete task successfully", nil))
}
