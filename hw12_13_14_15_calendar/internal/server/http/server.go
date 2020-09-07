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

// check requered fields
// compose event with coerce
// check fields validity

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
	r, ok := ctx.Value(repositoryKey).(repository.BaseRepo)

	if !ok {
		http.Error(w, "problem accessing DB", http.StatusInternalServerError)
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

func getEventFromReq(req *http.Request, userId repository.ID) (*repository.Event, error) {
	err := req.ParseForm()
	if err != nil {
		return nil, err
	}

	idStr := req.PostForm.Get("id")

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return nil, errors.New("failed to parse eventId")
	}

	// TODO: add validation for each field
	title := req.PostForm.Get("title")

	fmt.Printf("%+v", req.PostForm)
	//if title == "" {
	//	return nil, errors.New("specify title value")
	//}

	startAtStr := req.PostForm.Get("start_at")
	//if startAtStr == "" {
	//	return nil, errors.New("specify start_at value")
	//}

	start_at, err := getTimeFromTimestamp(startAtStr)
	if err != nil {
		return nil, err
	}

	endAtStr := req.PostForm.Get("end_at")
	//if endAtStr == "" {
	//	return nil, errors.New("specify end_at value")
	//}

	end_at, err := getTimeFromTimestamp(endAtStr)
	if err != nil {
		return nil, err
	}

	description := req.PostForm.Get("description")
	//if description == "" {
	//	return nil, errors.New("specify description value")
	//}

	notifyAtStr := req.PostForm.Get("notify_at")
	//if notifyAtStr == "" {
	//	return nil, errors.New("specify notify_at value")
	//}

	notify_at, err := getTimeFromTimestamp(notifyAtStr)
	if err != nil {
		return nil, err
	}

	return &repository.Event{
		ID:          repository.ID(id),
		Title:       title,
		StartAt:     start_at,
		EndAt:       end_at,
		Description: description,
		UserID:      userId,
		NotifyAt:    notify_at,
	}, nil
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
	r, ok := ctx.Value(repositoryKey).(repository.BaseRepo)

	if !ok {
		http.Error(w, "problem accessing DB", http.StatusInternalServerError)
		return
	}

	userId, err := getUserId(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event, err := getEventFromReq(req, userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = r.AddEvent(*event)

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
	r, ok := ctx.Value(repositoryKey).(repository.BaseRepo)

	if !ok {
		http.Error(w, "problem accessing DB", http.StatusInternalServerError)
		return
	}

	userId, err := getUserId(ctx)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	r, ok := ctx.Value(repositoryKey).(repository.BaseRepo)

	if !ok {
		http.Error(w, "problem accessing DB", http.StatusInternalServerError)
		return
	}

	userId, err := getUserId(ctx)

	// TODO: put inside middleware
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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

	// TODO: wrap log middleware on every handler
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
