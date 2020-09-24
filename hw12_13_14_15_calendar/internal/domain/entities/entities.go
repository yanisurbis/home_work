package entities

import "time"

type ID = int

type Event struct {
	ID          ID
	Title       string
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	Description string
	NotifyAt    time.Time `db:"notify_at"`
	UserID      int       `db:"user_id"`
}

type UpdateEventRequest struct {
	ID          ID
	Title       string
	StartAt     time.Time
	EndAt       time.Time
	Description string
	NotifyAt    time.Time
	UserID      int
}

type AddEventRequest struct {
	Title       string
	StartAt     time.Time
	EndAt       time.Time
	Description string
	NotifyAt    time.Time
	UserID      int
}

type DeleteEventRequest struct {
	ID     ID
	UserID ID
}

type GetEventsRequest struct {
	UserID ID
	Type   string
	From   time.Time
}

type User struct {
	ID
}
