package domain

import (
	"calendar/internal/domain/entities"
	"context"
	"time"
)

type EventStorage interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	AddEvent(event entities.Event) error
	UpdateEvent(userID entities.ID, event entities.Event) error
	DeleteEvent(eventID entities.ID) error
	GetEventsDay(userID entities.ID, from time.Time) ([]entities.Event, error)
	GetEventsWeek(userID entities.ID, from time.Time) ([]entities.Event, error)
	GetEventsMonth(userID entities.ID, from time.Time) ([]entities.Event, error)
	GetEventsToNotify(from time.Time, to time.Time) ([]entities.Event, error)
	DeleteOldEvents(to time.Time) error
	GetEvent(id entities.ID) (*entities.Event, error)
}
