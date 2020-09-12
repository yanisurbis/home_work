package http_server

import (
	"calendar/internal/repository"
	"context"
	"log"
	"net/http"
	"strconv"
	"time"
)

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			log.Println(r.RemoteAddr+" "+r.Method+" "+r.Host+" "+r.UserAgent(), " ", time.Since(t).Milliseconds(), "ms")
		}(time.Now())

		next.ServeHTTP(w, r)
	})
}

func createDbMiddleware(repo repository.BaseRepo) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, repositoryKey, repo)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func userIdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		userIdStr := r.Header.Get("userid")

		if userIdStr == "" {
			http.Error(w, "please specify userId in headers", http.StatusUnauthorized)
			return
		}

		userId, err := strconv.Atoi(userIdStr)

		if err != nil {
			http.Error(w, "please check your userId", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, userIdKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}