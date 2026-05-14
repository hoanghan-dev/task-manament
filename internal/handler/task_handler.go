package handler

import (
	"dev/task-management/internal/entities"
	"dev/task-management/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
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
	tasks := h.service.GetAllTask()
	c.JSON(200, tasks)
}

func (h *TaskHandler) GetTask(c *gin.Context) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		c.JSON(400, gin.H{
			"error": "Lỗi param nhé!",
		})
		return
	}
	task, err := h.service.GetTask(id)

	if err != nil {
		c.JSON(404, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(200, task)
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var taskReq entities.Task

	err := c.ShouldBindJSON(&taskReq)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	task := h.service.CreateTask(taskReq)
	c.JSON(201, task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var taskReq entities.Task

	err := c.ShouldBindJSON(&taskReq)

	idString := c.Param("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	task, err2 := h.service.UpdateTask(id, &taskReq)

	if err2 != nil {
		c.JSON(400, gin.H{
			"error": err2.Error(),
		})
		return
	}
	c.JSON(201, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idString := c.Param("id")

	id, err := strconv.Atoi(idString)

	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	err2 := h.service.DeleteTask(id)
	if err2 != nil {
		c.JSON(400, gin.H{
			"error": err2.Error(),
		})
		return
	}
	c.JSON(204, nil)
}
