package memory

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/domain/entities"
	"context"
	"fmt"
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

func getEvents(client *grpcclient.Client) []entities.Event {
	events, err := client.GetEventsDay(entities.GetEventsRequest{
		UserID: 1,
		Type:   "day",
		From:   time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
	return events
}

func testCreate(t *testing.T, client *grpcclient.Client) *entities.Event {
	addEventRequest := entities.AddEventRequest{
		Title:       "Event from test",
		StartAt:     time.Now().Add(3 * time.Hour),
		EndAt:       time.Now().Add(5 * time.Hour),
		Description: "Description from test",
		NotifyAt:    time.Now().Add(4 * time.Hour),
		UserID:      1,
	}


	err := client.AddEvent(addEventRequest)
	if err != nil {
		log.Fatal(err)
	}

	events := getEvents(client)

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

	return addedEvent
}

func testUpdate(t *testing.T, client *grpcclient.Client, addedEvent *entities.Event) *entities.Event {
	updateEventRequest := entities.UpdateEventRequest{
		ID:          addedEvent.ID,
		Title:       "new title",
		Description: "new description",
		UserID:      addedEvent.UserID,
	}

	err := client.UpdateEvent(updateEventRequest)
	if err != nil {
		log.Fatal(err)
	}

	events := getEvents(client)

	addedEvent1 := new(entities.Event)
	for _, event := range events {
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

	return addedEvent1
}

func testDelete(t *testing.T, client *grpcclient.Client, addedEvent *entities.Event) {
	deleteEventRequest := entities.DeleteEventRequest{
		ID:          addedEvent.ID,
		UserID:      addedEvent.UserID,
	}

	err := client.DeleteEvent(deleteEventRequest)
	if err != nil {
		log.Fatal(err)
	}

	events := getEvents(client)

	found := false
	for _, event := range events {
		if event.ID == deleteEventRequest.ID {
			found = true
		}
	}
	assert.False(t, found)
}

func testCRUD(t *testing.T, client *grpcclient.Client) {
	addedEvent := testCreate(t, client)
	_ = testUpdate(t, client, addedEvent)
	testDelete(t, client, addedEvent)
}

func testCRUDErrors(t *testing.T, client *grpcclient.Client) {
	requests := []entities.AddEventRequest{
		entities.AddEventRequest{
			Title:       "Event from test",
			StartAt:     time.Now().Add(-1 * time.Hour),
			EndAt:       time.Now().Add(5 * time.Hour),
			Description: "Description from test",
			NotifyAt:    time.Now().Add(4 * time.Hour),
			UserID:      1,
		},
	}

	responses := []error{}

	for _, request := range requests {
		responses = append(responses, client.AddEvent(request))
	}

	for _, response := range responses {
		fmt.Println(response)
		assert.Error(t, response)
	}
}

func TestIntegration(t *testing.T) {
	client := grpcclient.NewClient()
	err := client.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	t.Run("CRUD, basic cases work", func(t *testing.T) {
		testCRUD(t, client)
	})
	t.Run("CRUD, basic validations are present", func(t *testing.T) {
		testCRUDErrors(t, client)
	})
}
