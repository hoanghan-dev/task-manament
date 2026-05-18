package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestIdMiddleware(c *gin.Context) {
	requestId := uuid.New().String()

	c.Set("request_id", requestId)

	c.Header("X-request-ID", requestId)

	c.Next()
}
