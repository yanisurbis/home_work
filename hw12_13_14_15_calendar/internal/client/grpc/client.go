package grpcclient

import (
	"calendar/internal/domain/entities"
	"calendar/internal/lib"
	"calendar/internal/server/grpc/events_grpc"
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
)

type Client struct {
	conn   *grpc.ClientConn
	client events_grpc.EventsClient
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Start(ctx context.Context) error {
	conn, err := grpc.DialContext(ctx, "localhost:9090", grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = events_grpc.NewEventsClient(c.conn)

	return err
}

func convertEventsToNotifications(events []*events_grpc.EventResponse) ([]*entities.Notification, error) {
	notifications := make([]*entities.Notification, 0, len(events))

	for _, event := range events {
		startAt, err := lib.TimestampToTime(event.StartAt)
		if err != nil {
			return notifications, err
		}

		notification := entities.Notification{
			EventID:    entities.ID(event.Id),
			UserID:     entities.ID(event.UserId),
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
	log.Println(time.Now().Format(time.Stamp), "deleting old events")
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
