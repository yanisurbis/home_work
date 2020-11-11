package grpcclient

import (
	"calendar/internal/domain/entities"
	"calendar/internal/lib"
	"calendar/internal/server/grpc/events_grpc"
	"context"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/wrappers"

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
	// TODO: docker
	conn, err := grpc.DialContext(ctx, "localhost:9090", grpc.WithInsecure())
	if err != nil {
		return err
	}

	c.conn = conn
	c.client = events_grpc.NewEventsClient(c.conn)

	return err
}

func convertEventsResponseToEvents(response *events_grpc.EventsResponse) ([]entities.Event, error) {
	events := make([]entities.Event, 0, len(response.Events))

	for _, grpcEvent := range response.Events {
		startAt, err := lib.TimestampToTime(grpcEvent.StartAt)
		if err != nil {
			return nil, err
		}
		endAt, err := lib.TimestampToTime(grpcEvent.EndAt)
		if err != nil {
			return nil, err
		}
		notifyAt, err := lib.TimestampToTime(grpcEvent.NotifyAt)
		if err != nil {
			return nil, err
		}

		event := entities.Event{
			ID:          entities.ID(grpcEvent.Id),
			Title:       grpcEvent.Title,
			StartAt:     startAt,
			EndAt:       endAt,
			Description: grpcEvent.Description,
			NotifyAt:    notifyAt,
			UserID:      entities.ID(grpcEvent.UserId),
		}
		events = append(events, event)
	}

	return events, nil
}

func convertEventsToNotifications(
	events []*events_grpc.EventResponse,
) ([]*entities.Notification, error) {
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
		return nil, err
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

func (c *Client) AddEvent(request entities.AddEventRequest) error {
	startAt, err := ptypes.TimestampProto(request.StartAt)
	if err != nil {
		return err
	}
	endAt, err := ptypes.TimestampProto(request.EndAt)
	if err != nil {
		return err
	}
	notifyAt, err := ptypes.TimestampProto(request.NotifyAt)
	if err != nil {
		return err
	}

	r := events_grpc.AddEventRequest{
		Title:       request.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: request.Description,
		UserId:      uint32(request.UserID),
		NotifyAt:    notifyAt,
	}

	_, err = c.client.AddEvent(context.Background(), &r)
	return err
}

func (c *Client) UpdateEvent(request entities.UpdateEventRequest) error {
	startAt, err := ptypes.TimestampProto(request.StartAt)
	if err != nil {
		return err
	}
	endAt, err := ptypes.TimestampProto(request.EndAt)
	if err != nil {
		return err
	}
	notifyAt, err := ptypes.TimestampProto(request.NotifyAt)
	if err != nil {
		return err
	}

	r := events_grpc.UpdateEventRequest{
		Id:          uint32(request.ID),
		Title:       &wrappers.StringValue{Value: request.Title},
		StartAt:     startAt,
		EndAt:       endAt,
		Description: &wrappers.StringValue{Value: request.Description},
		UserId:      uint32(request.UserID),
		NotifyAt:    notifyAt,
	}

	_, err = c.client.UpdateEvent(context.Background(), &r)
	return err
}

func (c *Client) DeleteEvent(request entities.DeleteEventRequest) error {
	r := events_grpc.DeleteEventRequest{
		UserId:  uint32(request.UserID),
		EventId: uint32(request.ID),
	}

	_, err := c.client.DeleteEvent(context.Background(), &r)
	return err
}

func (c *Client) GetEventsDay(request entities.GetEventsRequest) ([]entities.Event, error) {
	from, err := ptypes.TimestampProto(request.From)
	if err != nil {
		return nil, err
	}

	r := events_grpc.GetEventsRequest{
		UserId: uint32(request.UserID),
		From:   from,
	}

	eventsResponse, err := c.client.GetEventsDay(context.Background(), &r)
	if err != nil {
		return nil, err
	}

	return convertEventsResponseToEvents(eventsResponse)

}

func (c *Client) GetEventsWeek(request entities.GetEventsRequest) ([]entities.Event, error) {
	from, err := ptypes.TimestampProto(request.From)
	if err != nil {
		return nil, err
	}

	r := events_grpc.GetEventsRequest{
		UserId: uint32(request.UserID),
		From:   from,
	}

	eventsResponse, err := c.client.GetEventsWeek(context.Background(), &r)
	if err != nil {
		return nil, err
	}

	return convertEventsResponseToEvents(eventsResponse)

}

func (c *Client) GetEventsMonth(request entities.GetEventsRequest) ([]entities.Event, error) {
	from, err := ptypes.TimestampProto(request.From)
	if err != nil {
		return nil, err
	}

	r := events_grpc.GetEventsRequest{
		UserId: uint32(request.UserID),
		From:   from,
	}

	eventsResponse, err := c.client.GetEventsMonth(context.Background(), &r)
	if err != nil {
		return nil, err
	}

	return convertEventsResponseToEvents(eventsResponse)

}

func (c *Client) Stop() error {
	return c.conn.Close()
}
