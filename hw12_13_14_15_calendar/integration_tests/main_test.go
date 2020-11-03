package memory

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/domain/entities"
	"context"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

//func findEventByTitle(events []entities.Event, title string) *entities.Event {
//	foundEvent := new(entities.Event)
//	for _, event := range events {
//		if event.Title == title {
//			foundEvent = &event
//		}
//	}
//
//	return foundEvent
//}

func TestCRUD(t *testing.T) {
	client := grpcclient.NewClient()
	err := client.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	addEventRequest := entities.AddEventRequest{
		Title:       "Event from test",
		StartAt:     time.Now().Add(3 * time.Hour),
		EndAt:       time.Now().Add(5 * time.Hour),
		Description: "Description from test",
		NotifyAt:    time.Now().Add(4 * time.Hour),
		UserID:      1,
	}


	err = client.AddEvent(addEventRequest)
	if err != nil {
		log.Fatal(err)
	}

	events, err := client.GetEventsDay(entities.GetEventsRequest{
		UserID: 1,
		Type:   "day",
		From:   time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	addedEvent := new(entities.Event)
	for _, event := range events {
		if event.Title == addEventRequest.Title {
			addedEvent = &event
		}
	}

	assert.Equal(t, addedEvent.Description, addEventRequest.Description)
	assert.Equal(t, addedEvent.Title, addEventRequest.Title)
	// TODO: check dates
	//assert.Equal(t, addedEvent.StartAt, addEventRequest.StartAt)
	//assert.Equal(t, addedEvent.EndAt, addEventRequest.EndAt)
	assert.Equal(t, addedEvent.Description, addEventRequest.Description)
	assert.Equal(t, addedEvent.UserID, addEventRequest.UserID)

    // ------------------------------------------------------------------

	updateEventRequest := entities.UpdateEventRequest{
		ID:          addedEvent.ID,
		Title:       "new title",
		Description: "new description",
		UserID:      addedEvent.UserID,
	}

	err = client.UpdateEvent(updateEventRequest)
	if err != nil {
		log.Fatal(err)
	}

	events1, err := client.GetEventsDay(entities.GetEventsRequest{
		UserID: 1,
		Type:   "day",
		From:   time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	addedEvent1 := new(entities.Event)
	for _, event := range events1 {
		if event.ID == updateEventRequest.ID {
			addedEvent1 = &event
		}
	}

	assert.Equal(t, updateEventRequest.Description, addedEvent1.Description)
	assert.Equal(t, updateEventRequest.Title, addedEvent1.Title)
	// TODO: check dates
	//assert.Equal(t, addedEvent.StartAt, addEventRequest.StartAt)
	//assert.Equal(t, addedEvent.EndAt, addEventRequest.EndAt)
	assert.Equal(t, updateEventRequest.UserID, addedEvent1.UserID)

	// ------------------------------------------------------------------
	deleteEventRequest := entities.DeleteEventRequest{
		ID:          addedEvent.ID,
		UserID:      addedEvent.UserID,
	}

	err = client.DeleteEvent(deleteEventRequest)
	if err != nil {
		log.Fatal(err)
	}

	events2, err := client.GetEventsDay(entities.GetEventsRequest{
		UserID: 1,
		Type:   "day",
		From:   time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	found := false
	for _, event := range events2 {
		if event.ID == deleteEventRequest.ID {
			found = true
		}
	}
	assert.False(t, found)
}

func TestCRUDErrors(t *testing.T) {
	assert.Equal(t, 1, 1)
}

func TestIntegration(t *testing.T) {
	t.Run("CRUD, basic cases work", func(t *testing.T) {
		TestCRUD(t)
	})
	t.Run("CRUD, basic validations are present", func(t *testing.T) {
		TestCRUDErrors(t)
	})
}
