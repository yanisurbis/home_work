package http_server

import (
	"calendar/internal/repository"
	"github.com/go-ozzo/ozzo-validation/v4"
)


//GET, events, query_params: from, type
//userId, from, type

//type Event struct {
//	ID          ID
//	Title       string
//	StartAt     time.Time `db:"start_at"`
//	EndAt       time.Time `db:"end_at"`
//	Description string
//	UserID      int       `db:"user_id"`
//	NotifyAt    time.Time `db:"notify_at"`
//}

//add-event
// should have user_id, handled by middleware
// required: all except id, notify_at, user_id
// check fields validity

//update-event
//delete-event

func validateEventToAdd(e repository.Event) error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&e.StartAt, validation.Required),
		validation.Field(&e.EndAt, validation.Required),
		validation.Field(&e.Description, validation.Required, validation.Length(1, 1000)),
		validation.Field(&e.UserID, validation.Required),
	)
}

func validateEventToUpdate(e repository.Event) error {
	return validation.ValidateStruct(&e,
		validation.Field(&e.Title, validation.Length(1, 100)),
		validation.Field(&e.Description, validation.Length(1, 1000)),
	)
}
