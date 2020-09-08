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

	from := time.Unix(int64(fromInt/1000), 0)
	fmt.Println("date -> " + from.String())

	return from, nil
}

func getFromParam(req *http.Request) (time.Time, error) {
	err := req.ParseForm()
	if err != nil {
		return time.Now(), err
	}

	fromStr := req.PostForm.Get("from")
	if fromStr == "" {
		return time.Now(), errors.New("specify from value")
	}

	return getTimeFromTimestamp(fromStr)
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

	events, err := cb(userId, from.Add(time.Duration(24)*time.Hour*-1), r)

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

func getEventToAdd(req *http.Request, userId repository.ID) (*repository.Event, error) {
	event := new(repository.Event)
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	title := req.PostForm.Get("title")
	event.Title = title

	if startAtStr := req.PostForm.Get("start_at"); startAtStr != "" {
		startAt, err := getTimeFromTimestamp(startAtStr)
		if err != nil {
			return nil, err
		}
		event.StartAt = startAt
	}

	if endAtStr := req.PostForm.Get("end_at"); endAtStr != "" {
		endAt, err := getTimeFromTimestamp(endAtStr)
		if err != nil {
			return nil, err
		}
		event.EndAt = endAt
	}

	description := req.PostForm.Get("description")
	event.Description = description

	if notifyAtStr := req.PostForm.Get("notify_at"); notifyAtStr != "" {
		notifyAt, err := getTimeFromTimestamp(notifyAtStr)
		if err != nil {
			return nil, err
		}
		event.NotifyAt = notifyAt
	}

	return event, nil
}

func getEventFromReqUpdate(req *http.Request, userId repository.ID, r repository.BaseRepo) (*repository.Event, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	idStr := req.PostForm.Get("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return nil, errors.New("failed to parse eventId")
	}

	event, err := r.GetEvent(userId, id)

	if err != nil {
		return nil, err
	}

	event.ID = id

	// TODO: add validation for each field
	title := req.PostForm.Get("title")
	if title != "" {
		event.Title = title
	}

	startAtStr := req.PostForm.Get("start_at")
	if startAtStr != "" {
		startAt, err := getTimeFromTimestamp(startAtStr)
		if err != nil {
			return nil, err
		}
		event.StartAt = startAt
	}

	endAtStr := req.PostForm.Get("end_at")
	if endAtStr != "" {
		endAt, err := getTimeFromTimestamp(endAtStr)
		if err != nil {
			return nil, err
		}
		event.EndAt = endAt
	}

	description := req.PostForm.Get("description")
	if description != "" {
		event.Description = description
	}

	notifyAtStr := req.PostForm.Get("notify_at")
	if notifyAtStr != "" {
		notifyAt, err := getTimeFromTimestamp(notifyAtStr)
		if err != nil {
			return nil, err
		}
		event.NotifyAt = notifyAt
	}

	return &event, nil
}

func addEvent(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	repo := getRepository(ctx)
	userId := getUserId(ctx)

	event, err := getEventToAdd(req, userId)
	if err != nil {
		// TODO: use standard error with field specification
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	w.Write(eventJSON)

}

func updateEvent(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	r := getRepository(ctx)
	userId := getUserId(ctx)

	event, err := getEventFromReqUpdate(req, userId, r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.UpdateEvent(userId, *event)

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

	http.HandleFunc("/hello", helloHandler)

	http.HandleFunc("/get-events-day", applyMiddlewares(getEventsDay, r))
	http.HandleFunc("/get-events-week", applyMiddlewares(getEventsWeek, r))
	http.HandleFunc("/get-events-month", applyMiddlewares(getEventsMonth, r))

	http.HandleFunc("/add-event", applyMiddlewares(addEvent, r))
	http.HandleFunc("/update-event", applyMiddlewares(updateEvent, r))
	http.HandleFunc("/delete-event", applyMiddlewares(deleteEvent, r))

	fmt.Println("server starting at port :8080")

	return s.instance.ListenAndServe()
}

func (s *Instance) Stop(ctx context.Context) error {
	return s.instance.Shutdown(ctx)
}
