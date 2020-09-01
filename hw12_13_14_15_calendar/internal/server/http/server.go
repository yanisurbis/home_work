package http_server

import (
	"calendar/internal/repository"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type Instance struct {
	instance *http.Server
}

const repositoryKey = "repository"

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

func applyMiddlewares(h BasicHandler, r repository.BaseRepo) BasicHandler {
	h1 := dbMiddleware(h, r)

	return logMiddleware(h1)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func getEvents(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	repo, ok := ctx.Value(repositoryKey).(repository.BaseRepo)

	if !ok {
		return
	}

	events, err := repo.GetEventsDay(1, time.Now().Add(time.Duration(24)*time.Hour*-1))

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
