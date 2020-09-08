package http_server

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
