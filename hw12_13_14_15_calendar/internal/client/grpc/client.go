package grpcclient

import (
	"calendar/internal/domain/entities"
	"calendar/internal/server/grpc/events_grpc"
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

func timestampToTime(ts *timestamppb.Timestamp) (time.Time, error) {
	if ts == nil {
		return time.Time{}, nil
	}

	return ptypes.Timestamp(ts)
}

func convertEventsToNotifications(events []*events_grpc.EventResponse) []*entities.Notification {
	// TODO: initialize length
	var notifications []*entities.Notification

	for _, event := range events {
		startAt, err := timestampToTime(event.StartAt)
		// TODO: handler errors
		if err != nil {
			return nil
		}

		notification := entities.Notification{
			EventId:    entities.ID(event.Id),
			UserId:     entities.ID(event.UserId),
			EventTitle: event.Title,
			StartAt:    startAt,
		}
		notifications = append(notifications, &notification)
	}

	return notifications
}

func GetNotifications() []*entities.Notification {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := events_grpc.NewEventsClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r := events_grpc.GetEventsRequest{
		UserId: 1,
		From:   &timestamp.Timestamp{Seconds: time.Now().Unix()},
	}

	res, err := c.GetEventsMonth(ctx, &r)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	return convertEventsToNotifications(res.Events)
}
