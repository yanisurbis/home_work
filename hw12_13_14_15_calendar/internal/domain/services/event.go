package domain

import (
	domain2 "calendar/internal/domain"
	"calendar/internal/domain/entities"
	"calendar/internal/domain/interfaces"
	"context"
	"github.com/go-ozzo/ozzo-validation/v4"
	"reflect"
	"time"
)

const (
	PeriodDay   = "day"
	PeriodWeek  = "week"
	PeriodMonth = "month"
)

type EventService struct {
	EventStorage domain.EventStorage
}

func validateEventToAdd(e entities.Event) error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&e.StartAt, validation.Required),
		validation.Field(&e.EndAt, validation.Required),
		validation.Field(&e.Description, validation.Required, validation.Length(1, 1000)),
		validation.Field(&e.UserID, validation.Required),
	)
}

func (es *EventService) AddEvent(ctx context.Context, title string, startAt time.Time, endAt time.Time,
									description string, notifyAt time.Time, userID entities.ID) (*entities.Event, error) {
	event := entities.Event{
		Title:       title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: description,
		NotifyAt:    notifyAt,
		UserID:      userID,
	}

	err := validateEventToAdd(event)

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

func validateEventToUpdate(e entities.Event) error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Title, validation.Length(1, 100)),
		validation.Field(&e.Description, validation.Length(1, 1000)),
	)
}

func mergeEvents(currEvent entities.Event, newEvent entities.Event) (*entities.Event, error) {
	err := validateEventToUpdate(newEvent)

	if err != nil {
		return nil, err
	}
	// TODO: we should check that startAt > endAt
	// TODO: we should check that startAt > curr

	if !reflect.ValueOf(newEvent.Title).IsZero() {
		currEvent.Title = newEvent.Title
	}
	if !reflect.ValueOf(newEvent.StartAt).IsZero() {
		currEvent.StartAt = newEvent.StartAt
	}
	if !reflect.ValueOf(newEvent.EndAt).IsZero() {
		currEvent.EndAt = newEvent.EndAt
	}
	// TODO: what should happen if user want to delete description?
	if !reflect.ValueOf(newEvent.Description).IsZero() {
		currEvent.Description = newEvent.Description
	}
	// TODO: what should happen if user want to delete notification
	if !reflect.ValueOf(newEvent.NotifyAt).IsZero() {
		currEvent.NotifyAt = newEvent.NotifyAt
	}

	return &currEvent, nil
}

//func (es *EventService) UpdateEvent(ctx context.Context, eventID entities.ID, title string, startAt time.Time, endAt time.Time,
//	description string, notifyAt time.Time, userID entities.ID) (*entities.Event, error) {
//	newEvent := entities.Event{
//		ID: 	     eventID,
//		Title:       title,
//		StartAt:     startAt,
//		EndAt:       endAt,
//		Description: description,
//		NotifyAt:    notifyAt,
//		UserID:      userID,
//	}
//
//	currEvent, err := es.EventStorage.GetEvent(userID, eventID)
//
//	if err != nil {
//		return nil, err
//	}
//
//	updatedEvent, err := mergeEvents(currEvent, newEvent)
//
//	if err != nil {
//		return nil, err
//	}
//
//	err = es.EventStorage.UpdateEvent(userID, *updatedEvent)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return updatedEvent, nil
//}

func (es *EventService) GetEvent(ctx context.Context, userID entities.ID, eventID entities.ID) (*entities.Event, error) {
	event, err := es.EventStorage.GetEvent(eventID)

	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, domain2.ErrNotFound
	}

	if event.UserID != userID {
		return nil, domain2.ErrForbidden
	}

	return event, nil
}

func (es *EventService) DeleteEvent(ctx context.Context, userID entities.ID, eventID entities.ID) (*entities.Event, error) {
	event, err := es.GetEvent(ctx, userID, eventID)

	if err != nil {
		return nil, err
	}

	err = es.EventStorage.DeleteEvent(userID, eventID)

	if err != nil {
		return nil, err
	}

	return event, nil
}

func (es *EventService) GetEvents(ctx context.Context, userID entities.ID, period string, from time.Time) ([]entities.Event, error) {
	if period == PeriodMonth {
		return es.EventStorage.GetEventsMonth(userID, from)
	} else if period == PeriodWeek {
		return es.EventStorage.GetEventsWeek(userID, from)
	} else if period == PeriodDay {
		return es.EventStorage.GetEventsDay(userID, from)
	}
	// TODO: log problem in case there is no match
	// TODO: make one storage function instead of 3
	return es.EventStorage.GetEventsDay(userID, from)
}