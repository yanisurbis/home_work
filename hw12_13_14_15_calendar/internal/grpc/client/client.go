package main

import (
	"calendar/internal/grpc/events_grpc"
	"context"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	client := events_grpc.NewEventsClient(conn)

	//from, err := ptypes.TimestampProto(time.Now().Add(time.Duration(20) * time.Hour * -1))



	/*from, err := ptypes.TimestampProto(time.Now().Add(time.Duration(20) * time.Hour * -1))

	if err != nil {
		log.Fatal(err)
	}

	query := &events_grpc.EventsQuery{
		UserId: 1,
		From:   from,
	}

	events, err := client.GetEventsDay(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", events)

	events, err = client.GetEventsWeek(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", events)

	events, err = client.GetEventsMonth(context.Background(), query)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v", events)*/
	userId := uint32(1)

	t, _ := ptypes.TimestampProto(time.Now())
	event := events_grpc.Event{
		Title:       "Event" + t.String(),
		StartAt:     t,
		EndAt:       t,
		Description: "Event" + t.String(),
		UserId:      userId,
		NotifyAt:    t,
	}

	_, err = client.AddEvent(context.Background(), &event)

	if err != nil {
		log.Fatal(err)
	}

	// -----------------------------------------------------

	from, _ := ptypes.TimestampProto(time.Now().Add(time.Duration(20) * time.Hour * -1))
	eventsResponse, err := client.GetEventsDay(context.Background(), &events_grpc.EventsQuery{
		UserId: userId,
		From:   from,
	})

	if eventsResponse == nil || len(eventsResponse.Events) != 1 {
		log.Fatal("adding event was unsuccessful")
	}

	// -----------------------------------------------------

	eventNew := event
	eventNew.Id = eventsResponse.Events[0].Id
	description := "new event"
	eventNew.Description = description

	_, err = client.UpdateEvent(context.Background(), &eventNew)

	if err != nil {
		log.Fatal(err)
	}

	eventsResponse, err = client.GetEventsDay(context.Background(), &events_grpc.EventsQuery{
		UserId: userId,
		From:   from,
	})

	if eventsResponse == nil || eventsResponse.Events[0].Description != description {
		log.Fatal("updating event was unsuccessful")
	}

	// -----------------------------------------------------

	_, err = client.DeleteEvent(context.Background(), &events_grpc.DeleteEventRequest{
		UserId:  userId,
		EventId: eventsResponse.Events[0].Id,
	})

	if err != nil {
		log.Fatal(err)
	}

	eventsResponse, err = client.GetEventsDay(context.Background(), &events_grpc.EventsQuery{
		UserId: userId,
		From:   from,
	})

	if eventsResponse == nil || len(eventsResponse.Events) != 0 {
		log.Fatal("deleting event was unsuccessful")
	}
}
