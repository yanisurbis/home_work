package sql

import (
	"calendar/internal/domain/entities"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/jmoiron/sqlx"
)

var (
	ErrForbidden = errors.New("access denied")
)

type Repo struct {
	db *sqlx.DB
}

func (r *Repo) Connect(ctx context.Context, dsn string) (err error) {
	ticker := backoff.NewTicker(backoff.NewExponentialBackOff())

	for range ticker.C {
		r.db, err = sqlx.Connect("pgx", dsn)
		if err != nil {
			log.Printf("could not connect to database: %+v", err)
			continue
		}
		log.Printf("connected successfully to database")
		return nil
	}

	return fmt.Errorf("failed to connect to database")
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) AddEvent(event entities.Event) (err error) {
	var events []entities.Event

	nstmt, err := r.db.PrepareNamed(
		"INSERT INTO events (title, start_at, end_at, description, user_id, notify_at) VALUES (:title, :start_at, :end_at, :description, :user_id, :notify_at)",
	)

	if err != nil {
		return
	}

	err = nstmt.Select(&events, event)

	return err
}

func (r *Repo) UpdateEvent(userID entities.ID, event entities.Event) (err error) {
	if userID != event.UserID {
		return ErrForbidden
	}

	var events []entities.Event

	nstmt, err := r.db.PrepareNamed(
		"UPDATE events SET title=:title, start_at=:start_at, end_at = :end_at, description = :description, notify_at=:notify_at WHERE  user_id = :user_id and id=:id",
	)

	if err != nil {
		return
	}

	err = nstmt.Select(&events, event)

	return
}

func (r *Repo) GetEvent(id entities.ID) (*entities.Event, error) {
	var events []entities.Event
	option := make(map[string]interface{})
	option["id"] = id

	nstmt, err := r.db.PrepareNamed("SELECT * FROM events WHERE id = :id")

	if err != nil {
		return nil, err
	}

	err = nstmt.Select(&events, option)

	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, nil
	}

	return &events[0], nil
}

func (r *Repo) DeleteEvent(eventID entities.ID) error {
	option := make(map[string]interface{})
	option["event_id"] = eventID

	nstmt, err := r.db.PrepareNamed("DELETE FROM events WHERE id=:event_id")

	if err != nil {
		return err
	}

	_, err = nstmt.Exec(option)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) getEvents(
	userID entities.ID,
	from time.Time,
	to time.Time,
) ([]entities.Event, error) {
	var events []entities.Event
	option := make(map[string]interface{})
	option["start"] = from
	option["end"] = to
	option["user_id"] = userID

	nstmt, err := r.db.PrepareNamed(
		"SELECT * FROM events WHERE user_id = :user_id and start_at>=:start and start_at<:end",
	)

	if err != nil {
		return nil, err
	}

	err = nstmt.Select(&events, option)

	if events == nil {
		return nil, nil
	}

	return events, err
}

func (r *Repo) GetEventsDay(userID entities.ID, from time.Time) ([]entities.Event, error) {
	return r.getEvents(userID, from, from.Add(time.Hour*24))
}

func (r *Repo) GetEventsWeek(userID entities.ID, from time.Time) ([]entities.Event, error) {
	return r.getEvents(userID, from, from.AddDate(0, 0, 7))
}

func (r *Repo) GetEventsMonth(userID entities.ID, from time.Time) ([]entities.Event, error) {
	return r.getEvents(userID, from, from.AddDate(0, 1, 0))
}

func (r *Repo) GetEventsToNotify(from time.Time, to time.Time) ([]entities.Event, error) {
	var events []entities.Event
	option := make(map[string]interface{})
	option["start"] = from
	option["end"] = to

	nstmt, err := r.db.PrepareNamed(
		"SELECT * FROM events WHERE notify_at >= :start and notify_at < :end",
	)
	if err != nil {
		return nil, err
	}

	err = nstmt.Select(&events, option)
	if err != nil {
		return nil, err
	}
	if events == nil {
		return nil, nil
	}

	return events, err
}

func (r *Repo) DeleteOldEvents(to time.Time) error {
	option := make(map[string]interface{})
	option["start"] = to

	nstmt, err := r.db.PrepareNamed("DELETE FROM events WHERE start_at < :start")

	if err != nil {
		return err
	}

	_, err = nstmt.Exec(option)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) AddNotifications(notifications []entities.Notification) error {
	nstmt, err := r.db.PrepareNamed(
		"INSERT INTO notifications (event_id, user_id, event_title, start_at) VALUES (:event_id, :user_id, :event_title, :start_at)",
	)

	if err != nil {
		return err
	}

	// TODO: Should use batch insert
	for _, notification := range notifications {
		_, err = nstmt.Exec(notification)
		if err != nil {
			return err
		}
	}

	return err
}

func (r *Repo) GetAllNotifications() ([]entities.Notification, error) {
	var notifications []entities.Notification

	nstmt, err := r.db.PrepareNamed("SELECT * FROM notifications ORDER BY event_id")
	if err != nil {
		return nil, err
	}

	option := make(map[string]interface{})
	err = nstmt.Select(&notifications, option)
	if err != nil {
		return nil, err
	}
	if notifications == nil {
		return nil, err
	}

	return notifications, err
}

func (r *Repo) DeleteAllNotifications() error {
	_, err := r.db.Exec("DELETE FROM notifications")

	if err != nil {
		return err
	}

	return nil
}
