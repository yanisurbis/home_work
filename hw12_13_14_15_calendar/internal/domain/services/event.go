package domain

import (
	"calendar/internal/domain/entities"
	"calendar/internal/domain/interfaces"
	"context"
	"time"
)

const (
	PERIOD_DAY = "day"
	PERIOD_WEEK = "week"
	PERIOD_MONTH = "month"
)

type EventService struct {
	EventStorage domain.EventStorage
}

func (es *EventService) AddEvent(ctx context.Context, title string, startAt time.Time, endAt time.Time,
									description string, notifyAt time.Time, userID entities.ID) (*entities.Event, error) {
	// TODO: add validation
	event := entities.Event{
		Title:       title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: description,
		NotifyAt:    notifyAt,
		UserID:      userID,
	}

	err := es.EventStorage.AddEvent(event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (es *EventService) UpdateEvent(ctx context.Context, eventID entities.ID, title string, startAt time.Time, endAt time.Time,
	description string, notifyAt time.Time, userID entities.ID) (*entities.Event, error) {
	event := entities.Event{
		ID: 	     eventID,
		Title:       title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: description,
		NotifyAt:    notifyAt,
		UserID:      userID,
	}

	event, err := es.EventStorage.GetEvent(userID, eventID)

	if err != nil {
		return nil, err
	}

	err = es.EventStorage.UpdateEvent(userID, event)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (es *EventService) DeleteEvent(ctx context.Context, userID entities.ID, eventID entities.ID) (*entities.Event, error) {
	event, err := es.EventStorage.GetEvent(userID, eventID)

	if err != nil {
		return nil, err
	}

	err = es.EventStorage.DeleteEvent(userID, eventID)

	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (es *EventService) GetEvents(ctx context.Context, userID entities.ID, period string, from time.Time) ([]entities.Event, error) {
	if period == PERIOD_MONTH {
		return es.EventStorage.GetEventsMonth(userID, from)
	} else if period == PERIOD_WEEK {
		return es.EventStorage.GetEventsWeek(userID, from)
	} else if period == PERIOD_DAY {
		return es.EventStorage.GetEventsDay(userID, from)
	}
	// TODO: log problem
	// TODO: make one storage function instead of 3
	return es.EventStorage.GetEventsDay(userID, from)
}