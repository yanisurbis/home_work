package memory

import (
	"calendar/internal/storage"
	"context"
	"errors"
	"sync"
	"time"
)

type DB struct {
	sync.Mutex
	events []storage.Event
}

func (m *DB) Connect(ctx context.Context, dsn string) error {
	return nil
}

func (m *DB) Close() error {
	return nil
}

func (m *DB) AddEvent(event storage.Event) error {
	m.Lock()
	m.events = append(m.events, event)
	m.Unlock()

	return nil
}

func (m *DB) UpdateEvent(event storage.Event) error {
	m.Lock()
	for i, e := range m.events {
		if e.ID == event.ID {
			if e.UserID != event.UserID {
				return errors.New("unauthorized request")
			}

			m.events[i] = event
		}
	}
	m.Unlock()

	return nil
}

func (m *DB) DeleteEvent(userID storage.ID, eventID storage.ID) error {
	var newEvents []storage.Event

	m.Lock()
	for _, e := range m.events {
		if e.ID == eventID {
			if e.UserID != userID {
				return errors.New("unauthorized request")
			}

			continue
		} else {
			newEvents = append(newEvents, e)
		}
	}

	m.events = newEvents
	m.Unlock()

	return nil
}

func filterDates(userID storage.ID, db *DB, from time.Time, to time.Time) []storage.Event {
	var dayEvents []storage.Event

	db.Lock()
	for _, e := range db.events {
		if e.UserID == userID && (e.StartAt.After(from) || e.StartAt.Equal(from)) && e.StartAt.Before(to) {
			dayEvents = append(dayEvents, e)
		}
	}
	db.Unlock()

	return dayEvents
}

func (m *DB) GetEventsDay(userID storage.ID, from time.Time) ([]storage.Event, error) {
	return filterDates(userID, m, from, from.AddDate(0, 0, 1)), nil
}

func (m *DB) GetEventsWeek(userID storage.ID, from time.Time) ([]storage.Event, error) {
	return filterDates(userID, m, from, from.AddDate(0, 0, 7)), nil
}

func (m *DB) GetEventsMonth(userID storage.ID, from time.Time) ([]storage.Event, error) {
	return filterDates(userID, m, from, from.AddDate(0, 1, 0)), nil
}
