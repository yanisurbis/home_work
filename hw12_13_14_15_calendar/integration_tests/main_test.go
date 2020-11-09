package memory

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/config"
	"calendar/internal/domain/entities"
	"calendar/internal/storage/sql"
	"context"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"strconv"
	"testing"
	"time"
)

func getEvents(client *grpcclient.Client, request entities.GetEventsRequest) []entities.Event {
	switch request.Type {
	case "day":
		events, err := client.GetEventsDay(request)
		if err != nil {
			log.Fatal(err)
		}
		return events
	case "week":
		events, err := client.GetEventsWeek(request)
		if err != nil {
			log.Fatal(err)
		}
		return events
	case "month":
		events, err := client.GetEventsMonth(request)
		if err != nil {
			log.Fatal(err)
		}
		return events
	default:
		events, err := client.GetEventsDay(request)
		if err != nil {
			log.Fatal(err)
		}
		return events
	}
}

func getEventsDay(client *grpcclient.Client) []entities.Event {
	return getEvents(client, entities.GetEventsRequest{
		UserID: 1,
		Type:   "day",
		From:   time.Now(),
	})
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

	events := getEventsDay(client)

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
	assert.Equal(t, 0, addedEvent.StartAt.Second()-addEventRequest.StartAt.Second())
	assert.Equal(t, 0, addedEvent.EndAt.Second()-addEventRequest.EndAt.Second())
	assert.Equal(t, 0, addedEvent.NotifyAt.Second()-addEventRequest.NotifyAt.Second())

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

	events := getEventsDay(client)

	updatedEvent := new(entities.Event)
	for _, event := range events {
		if event.ID == updateEventRequest.ID {
			updatedEvent = &event
		}
	}

	assert.Equal(t, entities.Event{
		ID:          updatedEvent.ID,
		Title:       updatedEvent.Title,
		Description: updatedEvent.Description,
		UserID:      updatedEvent.UserID,
	}, entities.Event{
		ID:          updateEventRequest.ID,
		Title:       updateEventRequest.Title,
		Description: updateEventRequest.Description,
		UserID:      updateEventRequest.UserID,
	})
	assert.Equal(t, 0, updatedEvent.StartAt.Second()-addedEvent.StartAt.Second())
	assert.Equal(t, 0, updatedEvent.EndAt.Second()-addedEvent.EndAt.Second())
	assert.Equal(t, 0, updatedEvent.NotifyAt.Second()-addedEvent.NotifyAt.Second())

	return updatedEvent
}

func testDelete(t *testing.T, client *grpcclient.Client, addedEvent *entities.Event) {
	deleteEventRequest := entities.DeleteEventRequest{
		ID:     addedEvent.ID,
		UserID: addedEvent.UserID,
	}

	err := client.DeleteEvent(deleteEventRequest)
	if err != nil {
		log.Fatal(err)
	}

	events := getEventsDay(client)

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

func checkEventsByIndex(t *testing.T, events []entities.Event, indexes []string) {
	var titles []string
	for _, event := range events {
		titles = append(titles, event.Title)
	}

	assert.Equal(t, indexes, titles)
}

func clearEvents(client *grpcclient.Client, events []entities.Event) {
	for _, event := range events {
		request := entities.DeleteEventRequest{
			ID:     event.ID,
			UserID: event.UserID,
		}
		err := client.DeleteEvent(request)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func testLists(t *testing.T, client *grpcclient.Client) {
	baseTime := time.Now().AddDate(1, 0, 0)
	baseUserId := 1

	baseEvent := entities.AddEventRequest{
		Title:       "Base event, title",
		StartAt:     baseTime,
		EndAt:       baseTime.AddDate(2, 0, 0),
		Description: "Base event, description",
		NotifyAt:    baseTime.Add(-2 * time.Hour),
		UserID:      baseUserId,
	}

	dates := []time.Time{
		baseTime.Add(1 * time.Hour),
		baseTime.Add(2 * time.Hour),
		baseTime.Add(3 * time.Hour),
		baseTime.AddDate(0, 0, 1),
		baseTime.AddDate(0, 0, 2),
		baseTime.AddDate(0, 0, 3),
		baseTime.AddDate(0, 0, 8),
		baseTime.AddDate(0, 0, 9),
		baseTime.AddDate(0, 0, 20),
		baseTime.AddDate(0, 2, 0),
		baseTime.AddDate(0, 3, 0),
		baseTime.AddDate(0, 4, 0),
	}

	for index, date := range dates {
		indexStr := strconv.Itoa(index)
		baseEvent.Title = indexStr
		baseEvent.StartAt = date
		err := client.AddEvent(baseEvent)
		if err != nil {
			log.Fatal(err)
		}
	}

	dayEvents := getEvents(client, entities.GetEventsRequest{
		UserID: baseUserId,
		Type:   "day",
		From:   baseTime,
	})
	checkEventsByIndex(t, dayEvents, []string{"0", "1", "2"})

	weekEvents := getEvents(client, entities.GetEventsRequest{
		UserID: baseUserId,
		Type:   "week",
		From:   baseTime,
	})
	checkEventsByIndex(t, weekEvents, []string{"0", "1", "2", "3", "4", "5"})

	monthEvents := getEvents(client, entities.GetEventsRequest{
		UserID: baseUserId,
		Type:   "month",
		From:   baseTime,
	})
	checkEventsByIndex(t, monthEvents, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8"})

	clearEvents(client, monthEvents)
}

func testEverything(t *testing.T, client *grpcclient.Client) {
	//TODO: delete when docker is set up
	err := os.Setenv("ENV", "TEST")
	if err != nil {
		log.Fatal(err)
	}

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	storage := new(sql.Repo)
	err = storage.Connect(context.Background(), c.PSQL.DSN)
	if err != nil {
		log.Fatal(err)
	}

	if err = storage.DeleteAllNotifications(); err != nil {
		log.Fatal(err)
	}

	baseTime := time.Now()
	requests := []entities.AddEventRequest{
		entities.AddEventRequest{
			Title:       "Test * Test * Test",
			StartAt:     baseTime.Add(1 * time.Minute),
			EndAt:       baseTime.Add(3 * time.Minute),
			Description: "Test event, description",
			NotifyAt:    baseTime.Add(1 * time.Second),
			UserID:      1,
		},
	}

	for _, request := range requests {
		err := client.AddEvent(request)
		if err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(15 * time.Second)

	dbNotifications, err := storage.GetAllNotifications()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Printf("%v\n", dbNotifications)
	for i, _ := range requests {
		r := requests[i]
		n := dbNotifications[i]
		assert.Equal(t, r.UserID, n.UserID)
		assert.Equal(t, r.Title, n.EventTitle)
		assert.Equal(t, 0, r.StartAt.Second()-n.StartAt.Second())
	}
}

func TestIntegration(t *testing.T) {
	client := grpcclient.NewClient()
	err := client.Start(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//t.Run("CRUD, basic cases work", func(t *testing.T) {
	//	testCRUD(t, client)
	//})
	//t.Run("CRUD, basic validations are present", func(t *testing.T) {
	//	testCRUDErrors(t, client)
	//})
	//t.Run("Check getEventsDay, getEventsWeek, getEventsMonth", func(t *testing.T) {
	//	testLists(t, client)
	//})
	t.Run("XXX", func(t *testing.T) {
		testEverything(t, client)
	})
}
