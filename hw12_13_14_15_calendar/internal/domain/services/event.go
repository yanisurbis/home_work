package domain

import (
	"calendar/internal/domain/entities"
	"calendar/internal/domain/errors"
	"calendar/internal/domain/interfaces"
	"context"
	"github.com/go-ozzo/ozzo-validation/v4"
	"time"
)

// TODO: move somewhere?
const (
	PeriodDay   = "day"
	PeriodWeek  = "week"
	PeriodMonth = "month"
)

var (
	// TODO: put random string here
	DefaultEmptyString = "_~_~_"
	DefaultEmptyTime   = time.Now().Add(-10)
)

type EventService struct {
	EventStorage domain.EventStorage
}

func validateEvent(e entities.Event) error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&e.StartAt, validation.Required),
		validation.Field(&e.EndAt, validation.Required),
		validation.Field(&e.Description, validation.Required, validation.Length(1, 1000)),
		validation.Field(&e.UserID, validation.Required),
	)
}

func (es *EventService) AddEvent(ctx context.Context, addEventRequest *entities.AddEventRequest) (*entities.Event, error) {
	event := entities.Event{
		Title:       addEventRequest.Title,
		StartAt:     addEventRequest.StartAt,
		EndAt:       addEventRequest.EndAt,
		Description: addEventRequest.Description,
		NotifyAt:    addEventRequest.NotifyAt,
		UserID:      addEventRequest.UserID,
	}

	err := validateEvent(event)

	if err != nil {
		return nil, err
	}

	// TODO: id should be created by us not to refetch db values
	err = es.EventStorage.AddEvent(event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func mergeEvents(currEvent *entities.Event, e *entities.UpdateEventRequest) (*entities.Event, error) {

	// TODO: we should check that startAt > endAt
	// TODO: we should check that startAt > curr

	if e.Title != DefaultEmptyString {
		currEvent.Title = e.Title
	}
	if !e.StartAt.IsZero() {
		currEvent.StartAt = e.StartAt
	}
	if !e.EndAt.IsZero() {
		currEvent.EndAt = e.EndAt
	}
	if e.Description != DefaultEmptyString {
		currEvent.Description = e.Description
	}
	if e.NotifyAt == DefaultEmptyTime {
		currEvent.NotifyAt = *new(time.Time)
	} else if !e.NotifyAt.IsZero() {
		currEvent.NotifyAt = e.NotifyAt
	}

	err := validateEvent(*currEvent)

	if err != nil {
		return nil, err
	}

	return currEvent, nil
}

func (es *EventService) UpdateEvent(ctx context.Context, eventUpdate *entities.UpdateEventRequest) (*entities.Event, error) {
	currEvent, err := es.GetEvent(ctx, eventUpdate.UserID, eventUpdate.ID)

	if err != nil {
		return nil, err
	}

	updatedEvent, err := mergeEvents(currEvent, eventUpdate)

	if err != nil {
		return nil, err
	}

	err = es.EventStorage.UpdateEvent(eventUpdate.UserID, *updatedEvent)

	if err != nil {
		return nil, err
	}

	return updatedEvent, nil
}

func (es *EventService) GetEvent(ctx context.Context, userID entities.ID, eventID entities.ID) (*entities.Event, error) {
	event, err := es.EventStorage.GetEvent(eventID)

	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.ErrNotFound
	}

	if event.UserID != userID {
		return nil, errors.ErrForbidden
	}

	return event, nil
}

func (es *EventService) DeleteEvent(ctx context.Context, deleteEventRequest *entities.DeleteEventRequest) (*entities.Event, error) {
	event, err := es.GetEvent(ctx, deleteEventRequest.UserID, deleteEventRequest.ID)

	if err != nil {
		return nil, err
	}

	err = es.EventStorage.DeleteEvent(deleteEventRequest.ID)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (es *EventService) GetEvents(ctx context.Context, getEventsRequest *entities.GetEventsRequest) ([]entities.Event, error) {
	period := getEventsRequest.Type
	userID := getEventsRequest.UserID
	from := getEventsRequest.From

	if period == PeriodMonth {
		return es.EventStorage.GetEventsMonth(userID, from)
	} else if period == PeriodWeek {
		return es.EventStorage.GetEventsWeek(userID, from)
	} else if period == PeriodDay {
		return es.EventStorage.GetEventsDay(userID, from)
	}
	// TODO: log problem in case there is no match
	return es.EventStorage.GetEventsDay(userID, from)
}
