package main

import (
	"calendar/internal/grpc/events_grpc"
	"context"
	"fmt"
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

	from, err := ptypes.TimestampProto(time.Now().Add(time.Duration(20) * time.Hour * -1))

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

	fmt.Printf("%+v", events)
}
