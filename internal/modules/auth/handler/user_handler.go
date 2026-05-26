package handler

import (
	"dev/task-management/internal/modules/auth/dto/request"
	"dev/task-management/internal/modules/auth/services"
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

	err := c.ShouldBindJSON(&userReq)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("validation failed", err.Error()))
		return
	}

	res, err := h.authService.Register(c.Request.Context(), &userReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("register account failed", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, response.ResponseSuccess("register account successfully", res))
}

func (h *AuthHander) Login(c *gin.Context) {
	var userLogin request.LoginUserRequestDTO
	err := c.ShouldBindJSON(&userLogin)

	if err != nil {
		c.JSON(http.StatusBadRequest, response.ResponseError("request format invalid", err.Error()))
		return
	}

	res, err := h.authService.Login(c.Request.Context(), &userLogin)

	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ResponseError("Authentication faild", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.ResponseSuccess("Authentication successfully", res))
}
