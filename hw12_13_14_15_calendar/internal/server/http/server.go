package http_server

import (
	"calendar/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Instance struct {
	instance *http.Server
}

const repositoryKey = "repository"
const userIdKey = "userId"

type BasicHandler func(http.ResponseWriter, *http.Request)

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
		userId := r.Header.Get("userid")
		ctx = context.WithValue(ctx, userIdKey, userId)

		h(w, r.WithContext(ctx))
	}
}

func applyMiddlewares(h BasicHandler, r repository.BaseRepo) BasicHandler {
	h1 := dbMiddleware(h, r)
	h2 := userIdMiddleware(h1)

	return logMiddleware(h2)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func getUserId(ctx context.Context) (repository.ID, error) {
	userId, ok := ctx.Value(userIdKey).(string)

	if !ok {
		return 0, errors.New("can't access userId")
	}

	if userId == "" {
		return 0, errors.New("specify userId in headers")
	}

	userIdInt, err := strconv.Atoi(userId)

	if err != nil {
		return 0, errors.New("can't convert userId")
	}

	return userIdInt, nil
}

func getEvents(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	repo, ok := ctx.Value(repositoryKey).(repository.BaseRepo)

	if !ok {
		http.Error(w, "problem accessing DB", http.StatusInternalServerError)
		return
	}

	userId, err := getUserId(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	events, err := repo.GetEventsDay(userId, time.Now().Add(time.Duration(24)*time.Hour*-1))

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", events)

	// TODO: handle empty array, right now return null
	// TODO: add migration to make file
	eventsJSON, err := json.Marshal(events)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(eventsJSON)
}

// TODO move grpc server in server folder
func (s *Instance) Start(r repository.BaseRepo) error {
	s.instance = &http.Server{Addr: ":8080"}

	// TODO: wrap log middleware on every handler
	http.HandleFunc("/get-events-day", applyMiddlewares(getEvents, r))
	fmt.Println("server starting at port :8080")

	return s.instance.ListenAndServe()
}

func (s *Instance) Stop(ctx context.Context) error {
	return s.instance.Shutdown(ctx)
}
