package postgres

import (
	"calendar/internal/domain"
	"calendar/internal/domain/entities"
	"calendar/internal/repository"
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrForbidden = errors.New("access denied")
)

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) Connect(ctx context.Context, dsn string) (err error) {
	r.db, err = sqlx.Connect("pgx", dsn)

	return
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) AddEvent(event entities.Event) (err error) {
	var events []entities.Event

	nstmt, err := r.db.PrepareNamed(
		"INSERT INTO events (title, start_at, end_at, description, user_id, notify_at) VALUES (:title, :start_at, :end_at, :description, :user_id, :notify_at)")

	if err != nil {
		return
	}

	err = nstmt.Select(&events, event)

	return err
}

func (r *Repo) UpdateEvent(userID repository.ID, event entities.Event) (err error) {
	if userID != event.UserID {
		return ErrForbidden
	}

	var events []entities.Event

	nstmt, err := r.db.PrepareNamed(
		"UPDATE events SET title=:title, start_at=:start_at, end_at = :end_at, description = :description, notify_at=:notify_at WHERE  user_id = :user_id and id=:id")

	if err != nil {
		return
	}

	err = nstmt.Select(&events, event)

	return
}

func (r *Repo) GetEvent(userId repository.ID, id repository.ID) (entities.Event, error) {
	var events []entities.Event
	option := make(map[string]interface{})
	option["id"] = id

	nstmt, err := r.db.PrepareNamed("SELECT * FROM events WHERE id = :id")

	if err != nil {
		return entities.Event{}, err
	}

	err = nstmt.Select(&events, option)

	if err != nil {
		// TODO: should we have not found error?
		// TODO: use pointers
		return entities.Event{}, nil
	}

	if len(events) == 0 {
		return entities.Event{}, domain.ErrNotFound
	}

	event := events[0]

	if event.UserID != userId {
		return entities.Event{}, ErrForbidden
	}

	return event, nil
}

func (r *Repo) DeleteEvent(userID repository.ID, eventID repository.ID) (err error) {
	var events []entities.Event
	option := make(map[string]interface{})
	option["event_id"] = eventID
	option["user_id"] = userID

	nstmt, err := r.db.PrepareNamed("DELETE FROM events WHERE id=:event_id and user_id = :user_id ")

	if err != nil {
		return
	}

	err = nstmt.Select(&events, option)

	return
}

func (r *Repo) getEvents(userID repository.ID, from time.Time, to time.Time) ([]entities.Event, error) {
	var events []entities.Event
	option := make(map[string]interface{})
	option["start"] = from
	option["end"] = to
	option["user_id"] = userID

	nstmt, err := r.db.PrepareNamed("SELECT * FROM events WHERE user_id = :user_id and start_at>=:start and start_at<:end")

	if err != nil {
		return nil, err
	}

	err = nstmt.Select(&events, option)

	return events, err
}

func (r *Repo) GetEventsDay(userID repository.ID, from time.Time) ([]entities.Event, error) {
	return r.getEvents(userID, from, from.Add(time.Hour*24))
}

func (r *Repo) GetEventsWeek(userID repository.ID, from time.Time) ([]entities.Event, error) {
	return r.getEvents(userID, from, from.AddDate(0, 0, 7))
}

func (r *Repo) GetEventsMonth(userID repository.ID, from time.Time) ([]entities.Event, error) {
	return r.getEvents(userID, from, from.AddDate(0, 1, 0))
}
