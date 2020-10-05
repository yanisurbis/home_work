package httpserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDStr := c.GetHeader("userId")
		userID, err := strconv.Atoi(userIDStr)

		if err != nil {
			c.String(http.StatusBadRequest, "please validate userId in headers")
			return
		}

		c.Set("userId", userID)
		c.Next()
	}
}