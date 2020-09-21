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
	UserID       int       `db:"user_id"`
}

type User struct {
	ID
}