package http_server

import (
	"calendar/internal/domain/entities"
	domain3 "calendar/internal/domain/errors"
	domain2 "calendar/internal/domain/interfaces"
	domain "calendar/internal/domain/services"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type Instance struct {
	instance *http.Server
}

const repositoryKey = "repository"
const userIdKey = "userId"

func getTimeFromTimestamp(timestamp string) (time.Time, error) {
	fromInt, err := strconv.Atoi(timestamp)
	if err != nil {
		return time.Now(), errors.New("can't convert from value")
	}

	from := time.Unix(int64(fromInt), 0)
	fmt.Println("date -> " + from.String())

	return from, nil
}

func prepareUpdateEventRequest(c *gin.Context) (*entities.UpdateEventRequest, error) {
	userId := GetUserID(c)

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		c.String(http.StatusBadRequest, "check eventId")
		return nil, err
	}

	eventUpdate := entities.UpdateEventRequest{}
	eventUpdate.ID = id
	eventUpdate.UserID = userId
	eventUpdate.Title = c.DefaultPostForm("title", domain.DefaultEmptyString)

	startAtStr := c.DefaultPostForm("start_at", domain.DefaultEmptyString)
	if startAtStr != domain.DefaultEmptyString {
		startAt, err := getTimeFromTimestamp(startAtStr)

		if err != nil {
			c.String(http.StatusBadRequest, "start_at wrong format")
			return nil, err
		}

		eventUpdate.StartAt = startAt
	}

	endAtStr := c.DefaultPostForm("end_at", domain.DefaultEmptyString)
	if endAtStr != domain.DefaultEmptyString {
		endAt, err := getTimeFromTimestamp(endAtStr)

		if err != nil {
			c.String(http.StatusBadRequest, "end_at wrong format")
			return nil, err
		}

		eventUpdate.EndAt = endAt
	}

	eventUpdate.Description = c.DefaultPostForm("description", domain.DefaultEmptyString)

	notifyAtStr := c.DefaultPostForm("notify_at", domain.DefaultEmptyString)
	if notifyAtStr != domain.DefaultEmptyString {
		if notifyAtStr == "" {
			eventUpdate.NotifyAt = domain.DefaultEmptyTime
		} else {
			notifyAt, err := getTimeFromTimestamp(notifyAtStr)

			if err != nil {
				c.String(http.StatusBadRequest, "notify_at wrong format")
				return nil, err
			}

			eventUpdate.NotifyAt = notifyAt
		}
	}

	return &eventUpdate, nil
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

		endAt, err := getTimeFromTimestamp(c.PostForm("end_at"))
		if err != nil {
			c.String(http.StatusBadRequest, "EndAt wrong format")
			return
		}

		description := c.PostForm("description")

		notifyAt, err := getTimeFromTimestamp(c.PostForm("notify_at"))
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

	router.PUT("/event/:id", func(c *gin.Context) {
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
	})

	fmt.Println("server starting at port :8080")

	return router.Run(":8080")
}

func (s *Instance) Stop(ctx context.Context) error {
	return s.instance.Shutdown(ctx)
}
