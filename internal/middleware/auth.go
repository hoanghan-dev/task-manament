package middleware

import (
	"dev/task-management/pkg/response"
	"dev/task-management/pkg/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	if header == "" {
		c.JSON(http.StatusUnauthorized, response.ResponseError("Unauthorized", "missing authorization header"))
		c.Abort()
		return
	}

	if !strings.HasPrefix(header, "Bearer ") {
		c.JSON(http.StatusUnauthorized, response.ResponseError("Unauthorized", "authorization header format invalid"))
		c.Abort()
		return
	}

	token := strings.TrimPrefix(header, "Bearer ")

	if token == "" {
		c.JSON(http.StatusUnauthorized, response.ResponseError("Unauthorized", "token is required"))
		c.Abort()
		return
	}

	payload, err := utils.VerifyAccessToken(token)

	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ResponseError("Unauthorized", err.Error()))
		c.Abort()
		return
	}

	email := payload["email"]

	emailStr, ok := email.(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, response.ResponseError("Unauthorized", "invalid token payload"))
		c.Abort()
		return
	}
	userId := payload["user_id"]
	userIdStr, ok := userId.(string)

	if !ok {
		c.JSON(http.StatusUnauthorized, response.ResponseError("Unauthorized", "invalid token payload"))
		c.Abort()
		return
	}

	c.Set("user_id", userIdStr)
	c.Set("email", emailStr)

	c.Next()
}
