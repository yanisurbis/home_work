package http_server

import (
	"calendar/internal/domain/entities"
	"github.com/gin-gonic/gin"
)

func GetUserID(c *gin.Context) entities.ID {
	userId, _ := c.Get("userId")
	return userId.(entities.ID)
}