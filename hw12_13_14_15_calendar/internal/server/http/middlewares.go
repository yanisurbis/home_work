package http_server

import (
	"calendar/internal/domain/entities"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)


func panicMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("recovered: ", err)
				http.Error(w, "Internal server error", 500)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			log.Println(r.RemoteAddr+" "+r.Method+" "+r.Host+" "+r.UserAgent(), " ", time.Since(t).Milliseconds(), "ms")
		}(time.Now())

		next.ServeHTTP(w, r)
	})
}

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

func GetUserID(c *gin.Context) entities.ID {
	userId, _ := c.Get("userId")
	return userId.(entities.ID)
}