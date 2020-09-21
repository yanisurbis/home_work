package http_server

import (
	domain3 "calendar/internal/domain/errors"
	domain2 "calendar/internal/domain/interfaces"
	domain "calendar/internal/domain/services"
	"calendar/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
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

func StatusOk(w http.ResponseWriter) {
	data := struct {
		Status string
	}{
		Status: "Ok",
	}

	dataJson, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(dataJson)

	if err != nil {
		log.Println(err)
	}
}

func StatusError(w http.ResponseWriter) {
	data := struct {
		Status string
	}{
		Status: "Ok",
	}

	dataJson, err := json.Marshal(data)

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(dataJson)

	if err != nil {
		log.Println(err)
	}
}

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

	StatusOk(w)
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

	StatusOk(w)
}

func getEventIdFromReq(req *http.Request) (repository.ID, error) {
	err := req.ParseForm()
	if err != nil {
		// TODO: should wrap the error
		return 0, err
	}

	eventIdStr := req.PostForm.Get("id")

	if eventIdStr == "" {
		return 0, validation.Errors{
			"Id": errors.New("event id is required"),
		}
	}

	eventId, err := strconv.Atoi(eventIdStr)

	if err != nil {
		return 0, validation.Errors{
			"Id": errors.New("wrong format"),
		}
	}

	return eventId, nil
}

func deleteEvent(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	r := getRepository(ctx)
	userId := getUserId(ctx)

	eventId, err := getEventIdFromReq(req)

	if err != nil {
		errs, _ := json.Marshal(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(errs)
		return
	}

	err = r.DeleteEvent(userId, eventId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	StatusOk(w)
}

// TODO move grpc server in server folder
func (s *Instance) Start1(r repository.BaseRepo) error {
	s.instance = &http.Server{Addr: ":8080"}

	router := mux.NewRouter()
	router.Use(panicMiddleware)
	router.Use(logMiddleware)
	router.HandleFunc("/hello", helloHandler)

	apiRouter := router.PathPrefix("/api").Subrouter()
	dbMiddleware := createDbMiddleware(r)
	apiRouter.Use(dbMiddleware, userIdMiddleware)

	apiRouter.HandleFunc("/events", getEventsMonth).Methods("GET").Queries("type", "month")
	apiRouter.HandleFunc("/events", getEventsWeek).Methods("GET").Queries("type", "week")
	apiRouter.HandleFunc("/events", getEventsDay).Methods("GET").Queries("type", "day")
	apiRouter.HandleFunc("/event", addEvent).Methods("POST")
	apiRouter.HandleFunc("/event", updateEvent).Methods("PUT")
	apiRouter.HandleFunc("/event", deleteEvent).Methods("PATCH")

	fmt.Println("server starting at port :8080")

	http.Handle("/", router)

	return s.instance.ListenAndServe()
}

// TODO: how to remove domain2?
func (s *Instance) Start(storage domain2.EventStorage) error {

	router := gin.Default()
	router.Use(UserIDMiddleware())
	// TODO: pass service, not storage
	eventService := domain.EventService{
		EventStorage: storage,
	}

	router.DELETE("/event/:id", func(c *gin.Context) {
		userId := GetUserID(c)
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)

		if err != nil {
			c.String(http.StatusBadRequest, "check eventId")
			return
		}

		deletedEvent, err := eventService.DeleteEvent(c, userId, id)

		if err != nil {
			if err == domain3.ErrForbidden {
				c.String(http.StatusForbidden, "don't have access")
				return
			} else if err == domain3.ErrNotFound {
				c.String(http.StatusNotFound, "event is not found")
				return
			} else {
				c.String(http.StatusInternalServerError, "error")
				return
			}
		}

		c.JSON(http.StatusOK, deletedEvent)
	})

	router.GET("/events", func(c *gin.Context) {
		userId := GetUserID(c)
		period := c.DefaultQuery("period", domain.PeriodDay)
		fromStr := c.Query("from")

		from, err := getTimeFromTimestamp(fromStr)

		if err != nil {
			c.String(http.StatusBadRequest, "check from parameter")
			return
		}

		events, err := eventService.GetEvents(c, userId, period, from)

		if err != nil {
			c.String(http.StatusInternalServerError, "error")
			return
		}

		c.JSON(http.StatusOK, events)
	})

	router.POST("/event", func(c *gin.Context) {
		userId := GetUserID(c)

		title := c.PostForm("title")

		startAt, err := getTimeFromTimestamp(c.PostForm("start_at"))
		if err != nil {
			c.String(http.StatusBadRequest, "StartAt wrong format")
			return
		}

		endAt, err := getTimeFromTimestamp(c.PostForm("start_at"))
		if err != nil {
			c.String(http.StatusBadRequest, "EndAt wrong format")
			return
		}

		description := c.PostForm("description")

		notifyAt, err := getTimeFromTimestamp(c.PostForm("start_at"))
		if err != nil {
			c.String(http.StatusBadRequest, "NotifyAt wrong format")
			return
		}

		addedEvent, err := eventService.AddEvent(c, title, startAt, endAt,
			description, notifyAt, userId)

		if err != nil {
			c.String(http.StatusInternalServerError, "error")
			return
		}

		c.JSON(http.StatusOK, addedEvent)
	})

	fmt.Println("server starting at port :8080")

	return router.Run(":8080")
}

func (s *Instance) Stop(ctx context.Context) error {
	return s.instance.Shutdown(ctx)
}
