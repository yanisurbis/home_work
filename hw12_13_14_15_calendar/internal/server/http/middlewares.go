package http_server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.GetHeader("userId")
		userId, err := strconv.Atoi(userIdStr)

		if err != nil {
			c.String(http.StatusBadRequest, "please validate userId in headers")
			return
		}

		c.Set("userId", userId)
		c.Next()
	}
}