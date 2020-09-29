package entities

import "time"

type ID = int

/*
db event
id uuid
title string, length > 0
start_at date
end_at date
description string, length >= 0
notify_at date, NULL if we don't have anything
user_id uuid
*/

type Event struct {
	ID          ID
	Title       string
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	Description string
	NotifyAt    time.Time `db:"notify_at"`
	UserID      ID       `db:"user_id"`
}

type UpdateEventRequest struct {
	ID          ID
	Title       string
	StartAt     time.Time
	EndAt       time.Time
	Description string
	NotifyAt    time.Time
	UserID      ID
}

type AddEventRequest struct {
	Title       string
	StartAt     time.Time
	EndAt       time.Time
	Description string
	NotifyAt    time.Time
	UserID      ID
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
