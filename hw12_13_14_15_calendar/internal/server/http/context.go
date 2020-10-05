package httpserver

import (
	"calendar/internal/domain/entities"

	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) entities.ID {
	userID, _ := c.Get("userId")

	return userID.(entities.ID)
}
