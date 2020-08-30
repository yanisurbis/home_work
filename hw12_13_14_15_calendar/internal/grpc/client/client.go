package main

import (
	"calendar/internal/grpc/events_grpc"
	"context"
	"google.golang.org/grpc"
	"log"
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

	//event := events_grpc.Event{
	//	Id:          7,
	//	Title:       "Updated event, aug30, 17:23",
	//	StartAt:     from,
	//	EndAt:       from,
	//	Description: "Updated event, aug30, 17:23",
	//	UserId:      1,
	//	NotifyAt:    from,
	//}

	//_, err = client.AddEvent(context.Background(), &event)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	//newEvent := event
	//
	//newEvent.Title = "Updated event"
	//
	//_, err = client.UpdateEvent(context.Background(), &newEvent)
	//
	//if err != nil {
	//	log.Fatal(err)
	//}

	deleteRequest := events_grpc.DeleteEventRequest{
		UserId:  1,
		EventId: 7,
	}

	_, err = client.DeleteEvent(context.Background(), &deleteRequest)

	if err != nil {
		log.Fatal(err)
	}
}
