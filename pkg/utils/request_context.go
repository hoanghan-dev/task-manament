package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetOwnerId(c *gin.Context) (uuid.UUID, error) {
	ownerIdStr := c.GetString("user_id")

	if ownerIdStr == "" {
		return uuid.UUID{}, errors.New("user_id not found in context after authentication")
	}

	ownerId, err := uuid.Parse(ownerIdStr)

	if err != nil {
		return uuid.UUID{}, errors.New("invalid user_id format in context")
	}

	return ownerId, nil
}
