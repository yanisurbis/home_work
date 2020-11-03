package memory

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/domain/entities"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCRUD(t *testing.T) {
	client := grpcclient.NewClient()
	err := client.Start(context.Background())
	if err != nil {
		fmt.Println(err.Error())
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
		fmt.Println(err.Error())
	}
	events, err := client.GetEventsDay(entities.GetEventsRequest{
		UserID: 1,
		Type:   "day",
		From:   time.Now(),
	})
	if err != nil {
		fmt.Println(err.Error())
	}

	addedEvent := new(entities.Event)
	for _, event := range events {
		if event.Title == addEventRequest.Title {
			addedEvent = &event
		}
	}

	assert.Equal(t, addedEvent.Description, addEventRequest.Description)
	assert.Equal(t, addedEvent.Title, addEventRequest.Title)
	//assert.Equal(t, addedEvent.StartAt, addEventRequest.StartAt)
	//assert.Equal(t, addedEvent.EndAt, addEventRequest.EndAt)
	assert.Equal(t, addedEvent.Description, addEventRequest.Description)
	assert.Equal(t, addedEvent.UserID, addEventRequest.UserID)
}

func TestIntegration(t *testing.T) {
	t.Run("random stuff", func(t *testing.T) {
		TestCRUD(t)
	})
}
