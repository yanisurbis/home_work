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

type Client struct {
	conn   *grpc.ClientConn
	client events_grpc.EventsClient
}

func NewClient() *Client {
	return &Client{}
}

func timestampToTime(ts *timestamppb.Timestamp) (time.Time, error) {
	if ts == nil {
		return time.Time{}, nil
	}

	return ptypes.Timestamp(ts)
}

func (c *Client) Start(cc context.Context) {
	ctx, _ := context.WithTimeout(cc, time.Second*15)
	conn, err := grpc.DialContext(ctx, "localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c.conn = conn
	c.client = events_grpc.NewEventsClient(c.conn)
}

func convertEventsToNotifications(events []*events_grpc.EventResponse) ([]*entities.Notification, error) {
	var notifications []*entities.Notification

	for _, event := range events {
		startAt, err := timestampToTime(event.StartAt)
		if err != nil {
			return notifications, err
		}

		notification := entities.Notification{
			EventId:    entities.ID(event.Id),
			UserId:     entities.ID(event.UserId),
			EventTitle: event.Title,
			StartAt:    startAt,
		}
		notifications = append(notifications, &notification)
	}

	return notifications, nil
}

func (c *Client) GetNotifications(from, to time.Time) ([]*entities.Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r := events_grpc.GetEventsToNotifyRequest{
		From: &timestamp.Timestamp{Seconds: from.Unix()},
		To:   &timestamp.Timestamp{Seconds: to.Unix()},
	}

	res, err := c.client.GetEventsToNotify(ctx, &r)
	if err != nil {
		return nil, nil
	}

	return convertEventsToNotifications(res.Events)
}

func (c *Client) DeleteOldEvents(to time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r := events_grpc.DeleteOldEventsRequest{
		To: &timestamp.Timestamp{Seconds: to.Unix()},
	}

	_, err := c.client.DeleteOldEvents(ctx, &r)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop() error {
	return c.conn.Close()
}
