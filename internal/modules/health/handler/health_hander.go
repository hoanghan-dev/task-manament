package health

import (
	"database/sql"
	"dev/task-management/internal/modules/health/dto"
	"dev/task-management/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler struct {
	redis    *redis.Client
	database *sql.DB
}

func NewHealthHandler(redis *redis.Client, database *sql.DB) *HealthHandler {
	return &HealthHandler{
		redis:    redis,
		database: database,
	}
}

func (h *HealthHandler) HealthCheck(c *gin.Context) {
	serverActive := true
	databaseActive := true
	redisActive := true

	err := h.database.Ping()
	if err != nil {
		databaseActive = false
	}

	err = h.redis.Ping(c.Request.Context()).Err()
	if err != nil {
		redisActive = false
	}

	healthRes := dto.NewHealthResponse(serverActive, redisActive, databaseActive)

	c.JSON(http.StatusOK, response.ResponseSuccess("server is running", healthRes))
}
