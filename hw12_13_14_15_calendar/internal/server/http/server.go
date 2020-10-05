package httpserver

import (
	"calendar/internal/domain/entities"
	domain3 "calendar/internal/domain/errors"
	domain "calendar/internal/domain/services"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Instance struct {
	instance *http.Server
}

func (s *Instance) Start(eventService domain.EventService) error {
	router := gin.Default()
	router.Use(UserIDMiddleware())

	router.DELETE("/event/:id", createDeleteEventHandler(eventService))
	router.GET("/events", createGetEventsHandler(eventService))
	router.POST("/event", createAddEventHandler(eventService))
	router.PUT("/event/:id", createUpdateEventHandler(eventService))

	fmt.Println("server starting at port :8080")
	s.instance = &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	if err := s.instance.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func createDeleteEventHandler(eventService domain.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		deleteEventRequest, err := prepareDeleteEventRequest(c)

		if err != nil {
			return
		}

		deletedEvent, err := eventService.DeleteEvent(c, deleteEventRequest)

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
	}
}

func createGetEventsHandler(eventService domain.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		getEventsRequest, err := prepareGetEventsRequest(c)

		if err != nil {
			return
		}

		events, err := eventService.GetEvents(c, getEventsRequest)

		if err != nil {
			c.String(http.StatusInternalServerError, "error")
			return
		}

		c.JSON(http.StatusOK, events)
	}
}

func createAddEventHandler(eventService domain.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		addEventRequest, err := prepareAddEventRequest(c)

		if err != nil {
			return
		}

		addedEvent, err := eventService.AddEvent(c, addEventRequest)

		if err != nil {
			c.String(http.StatusInternalServerError, "error")
			return
		}

		c.JSON(http.StatusOK, addedEvent)
	}
}

func createUpdateEventHandler(eventService domain.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		updateEventRequest, err := prepareUpdateEventRequest(c)

		if err != nil {
			return
		}

		addedEvent, err := eventService.UpdateEvent(c, updateEventRequest)

		if err != nil {
			c.String(http.StatusInternalServerError, "error")
			return
		}

		c.JSON(http.StatusOK, addedEvent)
	}
}

func getEventID(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("error converting event id")
	}

	return id, nil
}

func (s *Instance) Stop(ctx context.Context) error {
	return s.instance.Shutdown(ctx)
}

func timestampToTime(timestamp string) (time.Time, error) {
	fromInt, err := strconv.Atoi(timestamp)
	if err != nil {
		return time.Now(), errors.New("can't convert from value")
	}

	from := time.Unix(int64(fromInt), 0)

	return from, nil
}

func prepareDeleteEventRequest(c *gin.Context) (*entities.DeleteEventRequest, error) {
	deleteEventRequest := entities.DeleteEventRequest{}

	deleteEventRequest.UserID = GetUserID(c)

	id, err := getEventID(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return nil, err
	}
	deleteEventRequest.ID = id

	return &deleteEventRequest, nil
}

func prepareGetEventsRequest(c *gin.Context) (*entities.GetEventsRequest, error) {
	getEventsRequest := entities.GetEventsRequest{}

	getEventsRequest.UserID = GetUserID(c)
	getEventsRequest.Type = c.Query("period")

	fromStr := c.Query("from")
	from, err := timestampToTime(fromStr)
	if err != nil {
		c.String(http.StatusBadRequest, "check from parameter")
		return nil, err
	}
	getEventsRequest.From = from

	return &getEventsRequest, nil
}

func prepareAddEventRequest(c *gin.Context) (*entities.AddEventRequest, error) {
	addEventRequest := entities.AddEventRequest{}
	addEventRequest.UserID = GetUserID(c)

	addEventRequest.Title = c.PostForm("title")

	startAt, err := timestampToTime(c.PostForm("start_at"))
	if err != nil {
		c.String(http.StatusBadRequest, "StartAt wrong format")
		return nil, err
	}

	addEventRequest.StartAt = startAt

	endAt, err := timestampToTime(c.PostForm("end_at"))
	if err != nil {
		c.String(http.StatusBadRequest, "EndAt wrong format")
		return nil, err
	}

	addEventRequest.EndAt = endAt

	addEventRequest.Description = c.PostForm("description")

	notifyAtStr := c.PostForm("notify_at")
	if notifyAtStr != "" {
		notifyAt, err := timestampToTime(notifyAtStr)
		if err != nil {
			c.String(http.StatusBadRequest, "NotifyAt wrong format")
			return nil, err
		}
		addEventRequest.NotifyAt = notifyAt
	}

	return &addEventRequest, nil
}

func prepareUpdateEventRequest(c *gin.Context) (*entities.UpdateEventRequest, error) {
	userID := GetUserID(c)

	id, err := getEventID(c)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return nil, err
	}

	eventUpdate := entities.UpdateEventRequest{}
	eventUpdate.ID = id
	eventUpdate.UserID = userID

	eventUpdate.Title = c.PostForm("title")

	startAtStr := c.DefaultPostForm("start_at", domain.ValueNotPresent)
	if startAtStr != domain.ValueNotPresent {
		startAt, err := timestampToTime(startAtStr)

		if err != nil {
			c.String(http.StatusBadRequest, "start_at wrong format")
			return nil, err
		}

		eventUpdate.StartAt = startAt
	}

	endAtStr := c.DefaultPostForm("end_at", domain.ValueNotPresent)
	if endAtStr != domain.ValueNotPresent {
		endAt, err := timestampToTime(endAtStr)

		if err != nil {
			c.String(http.StatusBadRequest, "end_at wrong format")
			return nil, err
		}

		eventUpdate.EndAt = endAt
	}

	eventUpdate.Description = c.PostForm("description")

	notifyAtStr := c.DefaultPostForm("notify_at", domain.ValueNotPresent)
	if notifyAtStr != domain.ValueNotPresent {
		if notifyAtStr == "" {
			eventUpdate.NotifyAt = domain.ShouldResetTime
		} else {
			notifyAt, err := timestampToTime(notifyAtStr)

			if err != nil {
				c.String(http.StatusBadRequest, "notify_at wrong format")
				return nil, err
			}

			eventUpdate.NotifyAt = notifyAt
		}
	}

	return &eventUpdate, nil
}
