package storage

import (
	"context"
	"time"
)

type BaseRepo interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	AddEvent(event Event) error
	UpdateEvent(userID ID, event Event) error
	DeleteEvent(userID ID, eventID ID) error
	GetEventsDay(userID ID, from time.Time) ([]Event, error)
	GetEventsWeek(userID ID, from time.Time) ([]Event, error)
	GetEventsMonth(userID ID, from time.Time) ([]Event, error)
	GetEvent(userID ID, id ID) (Event, error)
}

type ID = int

type Event struct {
	ID          ID
	Title       string
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	Description string
	UserID      int       `db:"user_id"`
	NotifyAt    time.Time `db:"notify_at"`
}

type User struct {
	ID
}