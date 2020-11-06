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
	// TODO: remove add minute
	location, _ := time.LoadLocation("UTC")
	baseTime := time.Now().In(location).Add(1 * time.Minute)
	addEventRequest := entities.AddEventRequest{
		Title:       "Test event, title",
		StartAt:     baseTime,
		EndAt:       baseTime.Add(3 * time.Minute),
		Description: "Test event, description",
		NotifyAt:    baseTime.Add(-2 * time.Minute),
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

	assert.Equal(t, entities.Event{
		Title:       addedEvent.Title,
		Description: addedEvent.Description,
		UserID:      addedEvent.UserID,
	}, entities.Event{
		Title:       addEventRequest.Title,
		Description: addEventRequest.Description,
		UserID:      addEventRequest.UserID,
	})
	assert.Equal(t, 0, addedEvent.StartAt.Second() - addEventRequest.StartAt.Second())
	assert.Equal(t, 0, addedEvent.EndAt.Second() - addEventRequest.EndAt.Second())
	assert.Equal(t, 0, addedEvent.NotifyAt.Second() - addEventRequest.NotifyAt.Second())

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

	updatedEvent := new(entities.Event)
	for _, event := range events {
		if event.ID == updateEventRequest.ID {
			updatedEvent = &event
		}
	}

	assert.Equal(t, entities.Event{
		ID: updatedEvent.ID,
		Title:       updatedEvent.Title,
		Description: updatedEvent.Description,
		UserID:      updatedEvent.UserID,
	}, entities.Event{
		ID: updateEventRequest.ID,
		Title:       updateEventRequest.Title,
		Description: updateEventRequest.Description,
		UserID:      updateEventRequest.UserID,
	})
	assert.Equal(t, 0, updatedEvent.StartAt.Second() - addedEvent.StartAt.Second())
	assert.Equal(t, 0, updatedEvent.EndAt.Second() - addedEvent.EndAt.Second())
	assert.Equal(t, 0, updatedEvent.NotifyAt.Second() - addedEvent.NotifyAt.Second())

	return updatedEvent
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
	startAt := time.Now()
	endAt := startAt.Add(3 * time.Minute)
	notifyAt := startAt.Add(-2 * time.Minute)

	requests := []entities.AddEventRequest{
		// Notify after Start
		entities.AddEventRequest{
			Title:       "Event from test",
			StartAt:     startAt,
			EndAt:       endAt,
			Description: "Description from test",
			NotifyAt:    startAt.Add(2 * time.Minute),
			UserID:      1,
		},
		// Notify after Start
		entities.AddEventRequest{
			Title:       "Event from test",
			StartAt:     startAt,
			EndAt:       startAt.Add(-2 * time.Minute),
			Description: "Description from test",
			NotifyAt:    notifyAt,
			UserID:      1,
		},
		// Long title
		entities.AddEventRequest{
			Title:       "01234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890",
			StartAt:     startAt,
			EndAt:       endAt,
			Description: "Description from test",
			NotifyAt:    notifyAt,
			UserID:      1,
		},
	}

	responses := []error{}

	for _, request := range requests {
		responses = append(responses, client.AddEvent(request))
	}

	for _, response := range responses {
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
