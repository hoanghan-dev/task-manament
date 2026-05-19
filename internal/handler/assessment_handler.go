package handler

import (
	"dev/task-management/internal/dto/request"
	"dev/task-management/internal/services"
	"dev/task-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AssessmentHandler struct {
	AssessmentService services.AssessmentService
}

func NewAssessmentHandler(assessment services.AssessmentService) *AssessmentHandler {
	return &AssessmentHandler{
		AssessmentService: assessment,
	}
}

func (h *AssessmentHandler) GetAssessmentList(c *gin.Context) {
	assessments, err := h.AssessmentService.GetAssessmentList()
	if err != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("get assessment list failed", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.ResponseSuccess("get assessment list successfully", assessments))
}

func (h *AssessmentHandler) GetAssessment(c *gin.Context) {

	idString := c.Param("id")
	id, idErr := uuid.Parse(idString)

	if idErr != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", idErr.Error()))
		return
	}

	ass, err := h.AssessmentService.GetAssessment(id)

	if err != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("get assessment failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("get assessment successfully", ass))
}

func (h *AssessmentHandler) CreateAssessment(c *gin.Context) {
	var assCreate *request.AssessmentRequest

	assReqErr := c.ShouldBindJSON(&assCreate)

	if assReqErr != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", assReqErr.Error()))
		return
	}

	h.AssessmentService.CreateAssessment(assCreate)
	c.JSON(http.StatusCreated, response.ResponseSuccess("create assessment successfully", nil))
}

func (h *AssessmentHandler) UpdateAssessment(c *gin.Context) {
	var assCreate *request.AssessmentRequest
	idString := c.Param("id")
	id, idErr := uuid.Parse(idString)
	assReqErr := c.ShouldBindJSON(&assCreate)

	if idErr != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", idErr.Error()))
		return
	}

	if assReqErr != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", assReqErr.Error()))
		return
	}

	err := h.AssessmentService.UpdateAssessment(id, assCreate)

	if err != nil {
		c.JSON(http.StatusNotFound, response.ResponseError("update assessment failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("update assessment successfully", nil))
}

func (h *AssessmentHandler) DeleteAssessment(c *gin.Context) {
	id, idErr := uuid.Parse(c.Param("id"))

	if idErr != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", idErr.Error()))
		return
	}

	err := h.AssessmentService.DeleteAssessment(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("delete assessment failed", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("delete assessment successfully", nil))
}
