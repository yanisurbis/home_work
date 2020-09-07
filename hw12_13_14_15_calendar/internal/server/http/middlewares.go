package http_server

import (
	"calendar/internal/repository"
	"log"
	"net/http"
	"strconv"
	"time"
	"context"
)

func logMiddleware(h BasicHandler) BasicHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(t time.Time) {
			log.Println(r.RemoteAddr+" "+r.Method+" "+r.Host+" "+r.UserAgent(), " ", time.Since(t).Milliseconds(), "ms")
		}(time.Now())

		h(w, r)
	}
}

func dbMiddleware(h BasicHandler, repo repository.BaseRepo) BasicHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, repositoryKey, repo)

		h(w, r.WithContext(ctx))
	}
}

func userIdMiddleware(h BasicHandler) BasicHandler {
	return func(w http.ResponseWriter, r *http.Request) {
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
		h(w, r.WithContext(ctx))
	}
}

func applyMiddlewares(h BasicHandler, r repository.BaseRepo) BasicHandler {
	h1 := dbMiddleware(h, r)
	h2 := userIdMiddleware(h1)

	return logMiddleware(h2)
}