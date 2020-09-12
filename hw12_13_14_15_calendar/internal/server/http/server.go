package http_server

import (
	"calendar/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/mux"
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

// check requered fields
// compose event with coerce
// check fields validity

func helloHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello world\n")
}

func getTimeFromTimestamp(timestamp string) (time.Time, error) {
	fromInt, err := strconv.Atoi(timestamp)
	if err != nil {
		return time.Now(), errors.New("can't convert from value")
	}

	from := time.Unix(int64(fromInt), 0)
	fmt.Println("date -> " + from.String())

	return from, nil
}

func getFromParam(req *http.Request) (time.Time, error) {
	fromValues, ok := req.URL.Query()["from"]

	if !ok || len(fromValues) == 0 {
		return time.Now(), errors.New("specify from value")
	}

	return getTimeFromTimestamp(fromValues[0])
}

func getEvents(w http.ResponseWriter, req *http.Request, cb func(userID repository.ID, from time.Time, repo repository.BaseRepo) ([]repository.Event, error)) {
	ctx := req.Context()

	r := getRepository(ctx)
	userId := getUserId(ctx)

	from, err := getFromParam(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: remove ADD
	events, err := cb(userId, from.Add(time.Duration(23)*time.Hour*-1), r)

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

func getEventsDay(w http.ResponseWriter, req *http.Request) {
	getEvents(w, req, func(userId repository.ID, from time.Time, repo repository.BaseRepo) ([]repository.Event, error) {
		return repo.GetEventsDay(userId, from)
	})
}

func getEventsWeek(w http.ResponseWriter, req *http.Request) {
	getEvents(w, req, func(userId repository.ID, from time.Time, repo repository.BaseRepo) ([]repository.Event, error) {
		return repo.GetEventsWeek(userId, from)
	})
}

func getEventsMonth(w http.ResponseWriter, req *http.Request) {
	getEvents(w, req, func(userId repository.ID, from time.Time, repo repository.BaseRepo) ([]repository.Event, error) {
		return repo.GetEventsMonth(userId, from)
	})
}

func parseEventToAdd(req *http.Request, userId repository.ID) (*repository.Event, error) {
	event := new(repository.Event)
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	event.UserID = userId

	title := req.PostForm.Get("title")
	event.Title = title

	if startAtStr := req.PostForm.Get("start_at"); startAtStr != "" {
		startAt, err := getTimeFromTimestamp(startAtStr)
		if err != nil {
			return nil, validation.Errors{
				"StartAt": errors.New("wrong format"),
			}
		}
		event.StartAt = startAt
	}

	if endAtStr := req.PostForm.Get("end_at"); endAtStr != "" {
		endAt, err := getTimeFromTimestamp(endAtStr)
		if err != nil {
			return nil, validation.Errors{
				"EndAt": errors.New("wrong format"),
			}
		}
		event.EndAt = endAt
	}

	description := req.PostForm.Get("description")
	event.Description = description

	if notifyAtStr := req.PostForm.Get("notify_at"); notifyAtStr != "" {
		notifyAt, err := getTimeFromTimestamp(notifyAtStr)
		if err != nil {
			return nil, validation.Errors{
				"NotifyAt": errors.New("wrong format"),
			}
		}
		event.NotifyAt = notifyAt
	}

	if err = validateEventToAdd(*event); err != nil {
		return nil, err
	}

	return event, nil
}

func parseEventToUpdate(req *http.Request, userId repository.ID, r repository.BaseRepo) (*repository.Event, error) {
	err := req.ParseForm()
	if err != nil {
		// TODO: consistent format for error?
		return nil, err
	}

	idStr := req.PostForm.Get("id")
	id := 0
	if idStr == "" {
		return nil, validation.Errors{
			"Id": errors.New("event id is required"),
		}
	} else {
		id, err = strconv.Atoi(idStr)

		if err != nil {
			return nil, validation.Errors{
				"Id": errors.New("wrong format"),
			}
		}
	}

	event, err := r.GetEvent(userId, id)
	if err != nil {
		// TODO: consistent format for error?
		// TODO: handle panic
		return nil, err
	}

	if title := req.PostForm.Get("title"); title != "" {
		event.Title = title
	}

	if startAtStr := req.PostForm.Get("start_at"); startAtStr != "" {
		startAt, err := getTimeFromTimestamp(startAtStr)
		if err != nil {
			return nil, validation.Errors{
				"StartAt": errors.New("wrong format"),
			}
		}
		event.StartAt = startAt
	}

	if endAtStr := req.PostForm.Get("end_at"); endAtStr != "" {
		endAt, err := getTimeFromTimestamp(endAtStr)
		if err != nil {
			return nil, validation.Errors{
				"EndAt": errors.New("wrong format"),
			}
		}
		event.EndAt = endAt
	}

	if description := req.PostForm.Get("description"); description != "" {
		event.Description = description
	}

	if notifyAtStr := req.PostForm.Get("notify_at"); notifyAtStr != "" {
		notifyAt, err := getTimeFromTimestamp(notifyAtStr)
		if err != nil {
			return nil, validation.Errors{
				"NotifyAt": errors.New("wrong format"),
			}
		}
		event.NotifyAt = notifyAt
	}

	if err = validateEventToUpdate(event); err != nil {
		return nil, err
	}

	return &event, nil
}

func addEvent(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	repo := getRepository(ctx)
	userId := getUserId(ctx)

	event, err := parseEventToAdd(req, userId)
	if err != nil {
		fmt.Printf("%T", err)
		errs, _ := json.Marshal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errs)
		return
	}

	err = repo.AddEvent(*event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// send status, not event
	w.Header().Set("Content-Type", "application/json")
	w.Write(eventJSON)
}

func updateEvent(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	repo := getRepository(ctx)
	userId := getUserId(ctx)

	event, err := parseEventToUpdate(req, userId, repo)

	if err != nil {
		errs, _ := json.Marshal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errs)
		return
	}

	err = repo.UpdateEvent(userId, *event)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Ok"))
}

func getEventIdFromReq(req *http.Request) (repository.ID, error) {
	err := req.ParseForm()
	if err != nil {
		// TODO: should wrap the error
		return 0, err
	}

	// TODO: change to id instead of eventId
	eventIdStr := req.PostForm.Get("eventId")

	if eventIdStr == "" {
		return 0, errors.New("specify eventId value")
	}

	eventId, err := strconv.Atoi(eventIdStr)

	if err != nil {
		return 0, errors.New("failed to parse eventId")
	}

	return eventId, nil
}

func deleteEvent(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	r := getRepository(ctx)
	userId := getUserId(ctx)

	eventId, err := getEventIdFromReq(req)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.DeleteEvent(userId, eventId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Ok"))

	if err != nil {
		log.Println("Failed to send response")
		return
	}
}

// TODO move grpc server in server folder
func (s *Instance) Start(r repository.BaseRepo) error {
	s.instance = &http.Server{Addr: ":8080"}

	router := mux.NewRouter()

	dbMiddleware := createDbMiddleware(r)
	router.Use(logMiddleware, dbMiddleware, userIdMiddleware)

	router.HandleFunc("/hello", helloHandler)
	router.HandleFunc("/events", getEventsMonth).Methods("GET").Queries("type", "month")
	router.HandleFunc("/events", getEventsWeek).Methods("GET").Queries("type", "week")
	router.HandleFunc("/events", getEventsDay).Methods("GET").Queries("type", "day")
	router.HandleFunc("/event", addEvent).Methods("POST")
	router.HandleFunc("/event", updateEvent).Methods("PUT")
	router.HandleFunc("/event", deleteEvent).Methods("DELETE")

	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.HandleFunc("/hello-world", helloHandler)

	fmt.Println("server starting at port :8080")

	http.Handle("/", router)

	return s.instance.ListenAndServe()
}

func (s *Instance) Stop(ctx context.Context) error {
	return s.instance.Shutdown(ctx)
}
