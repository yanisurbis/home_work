package domain

import (
	"calendar/internal/domain/entities"
	domainErrors "calendar/internal/domain/errors"
	domain "calendar/internal/domain/interfaces"
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	PeriodDay   = "day"
	PeriodWeek  = "week"
	PeriodMonth = "month"
)

var (
	ShouldResetString = "should_reset_string"
	ValueNotPresent   = "value_not_present"
	ShouldResetTime   = time.Now().Add(-10)
)

type EventService struct {
	EventStorage domain.EventStorage
}

func validateEvent(e entities.Event) error {
	if e.StartAt.Before(time.Now()) {
		return errors.New("start_at should be grater than current date")
	}

	if e.EndAt.Before(e.StartAt) {
		return errors.New("end_at should not be less than start_at")
	}

	if !e.NotifyAt.IsZero() && e.NotifyAt.Before(e.StartAt) {
		return errors.New("notify_at should not be less than start_at")
	}

	return validation.ValidateStruct(&e,
		validation.Field(&e.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&e.StartAt, validation.Required),
		validation.Field(&e.EndAt, validation.Required),
		validation.Field(&e.Description, validation.Length(0, 1000)),
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

	err = es.EventStorage.AddEvent(event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func mergeEvents(currEvent *entities.Event, e *entities.UpdateEventRequest) (*entities.Event, error) {
	if !e.StartAt.IsZero() {
		currEvent.StartAt = e.StartAt
	}
	if !e.EndAt.IsZero() {
		currEvent.EndAt = e.EndAt
	}
	if e.Description != ShouldResetString {
		currEvent.Description = e.Description
	}
	if e.NotifyAt == ShouldResetTime {
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
		return nil, domainErrors.ErrNotFound
	}

	if event.UserID != userID {
		return nil, domainErrors.ErrForbidden
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

	switch period {
	case PeriodMonth:
		return es.EventStorage.GetEventsMonth(userID, from)
	case PeriodWeek:
		return es.EventStorage.GetEventsWeek(userID, from)
	case PeriodDay:
		return es.EventStorage.GetEventsDay(userID, from)
	default:
		return []entities.Event{}, nil
	}
}

func (es *EventService) GetEventsToNotify(ctx context.Context, getEventsRequest *entities.GetEventsToNotifyRequest) ([]entities.Event, error) {
	return es.EventStorage.GetEventsToNotify(getEventsRequest.UserID, getEventsRequest.From, getEventsRequest.To)
}
