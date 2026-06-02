package handler

import (
	"dev/task-management/internal/modules/auth/dto/request"
	"dev/task-management/internal/modules/auth/services"
	"dev/task-management/pkg/apperror"
	"dev/task-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHander struct {
	authService services.AuthService
}

func NewAuthHandler(authS services.AuthService) *AuthHander {
	return &AuthHander{
		authService: authS,
	}
}

func (h *AuthHander) Register(c *gin.Context) {
	var userReq request.RegisterUserRequestDTO

	if err := c.ShouldBindJSON(&userReq); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	res, err := h.authService.Register(c.Request.Context(), &userReq)
	if err != nil {
		// 400 Validation / 409 Conflict (duplicate email) / 500 Internal
		apperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response.ResponseSuccess("register account successfully", res))
}

func (h *AuthHander) Login(c *gin.Context) {
	var userLogin request.LoginUserRequestDTO

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		apperror.HandleError(c, apperror.NewValidation(err.Error()))
		return
	}

	res, err := h.authService.Login(c.Request.Context(), &userLogin)
	if err != nil {
		// 401 Unauthorized / 500 Internal
		apperror.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("authentication successfully", res))
}
