package servergrpc

import (
	"calendar/internal/domain/entities"
	domain "calendar/internal/domain/services"
	"calendar/internal/storage"
	"calendar/internal/server/grpc/events_grpc"
	"context"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	eventService domain.EventService
	instance     *grpc.Server
}

func timestampToTime(ts *timestamppb.Timestamp) (time.Time, error) {
	if ts == nil {
		return time.Time{}, nil
	}

	return ptypes.Timestamp(ts)
}

func createEventResponse(event entities.Event) (*events_grpc.EventResponse, error) {
	startAt, err := ptypes.TimestampProto(event.StartAt)
	if err != nil {
		return nil, errors.New("start_at conversion error")
	}

	endAt, err := ptypes.TimestampProto(event.EndAt)
	if err != nil {
		return nil, errors.New("end_at conversion error")
	}

	var notifyAt *timestamp.Timestamp = nil
	if !event.NotifyAt.IsZero() {
		notifyAt, err = ptypes.TimestampProto(event.NotifyAt)
		if err != nil {
			return nil, errors.New("notify_at conversion error")
		}
	}

	return &events_grpc.EventResponse{
		Id:          uint32(event.ID),
		Title:       event.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: event.Description,
		UserId:      uint32(event.UserID),
		NotifyAt:    notifyAt,
	}, nil
}

func (s *Server) GetEvents(ctx context.Context, query *events_grpc.GetEventsRequest, period string) (*events_grpc.EventsResponse, error) {
	from, err := timestampToTime(query.From)
	if err != nil {
		return nil, errors.Wrap(err, "from field conversion error")
	}

	getEventsRequest := entities.GetEventsRequest{
		UserID: storage.ID(query.UserId),
		Type:   period,
		From:   from,
	}
	events, err := s.eventService.GetEvents(ctx, &getEventsRequest)
	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch events")
	}

	eventsResponse := []*events_grpc.EventResponse{}
	for _, event := range events {
		event, err := createEventResponse(event)
		if err != nil {
			return nil, err
		}
		eventsResponse = append(eventsResponse, event)
	}

	return &events_grpc.EventsResponse{Events: eventsResponse}, nil
}

func (s *Server) GetEventsDay(ctx context.Context, query *events_grpc.GetEventsRequest) (*events_grpc.EventsResponse, error) {
	return s.GetEvents(ctx, query, domain.PeriodDay)
}

func (s *Server) GetEventsWeek(ctx context.Context, query *events_grpc.GetEventsRequest) (*events_grpc.EventsResponse, error) {
	return s.GetEvents(ctx, query, domain.PeriodWeek)
}

func (s *Server) GetEventsMonth(ctx context.Context, query *events_grpc.GetEventsRequest) (*events_grpc.EventsResponse, error) {
	return s.GetEvents(ctx, query, domain.PeriodMonth)
}

func prepareAddEventRequest(eventGrpc *events_grpc.AddEventRequest) (*entities.AddEventRequest, error) {
	startAt, err := timestampToTime(eventGrpc.StartAt)
	if err != nil {
		return nil, errors.New("error converting event.startAt")
	}

	endAt, err := timestampToTime(eventGrpc.EndAt)
	if err != nil {
		return nil, errors.New("error converting event.endAt")
	}

	notifyAt, err := timestampToTime(eventGrpc.NotifyAt)
	if err != nil {
		return nil, errors.New("error converting event.notifyAt")
	}

	return &entities.AddEventRequest{
		Title:       eventGrpc.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: eventGrpc.Description,
		NotifyAt:    notifyAt,
		UserID:      storage.ID(eventGrpc.UserId),
	}, nil
}

func (s *Server) AddEvent(ctx context.Context, query *events_grpc.AddEventRequest) (*events_grpc.EventResponse, error) {
	addEventRequest, err := prepareAddEventRequest(query)

	if err != nil {
		return nil, err
	}

	event, err := s.eventService.AddEvent(ctx, addEventRequest)

	if err != nil {
		return nil, errors.Wrap(err, "problem adding event to the DB")
	}

	return createEventResponse(*event)
}

func prepareUpdateEventRequest(eventGrpc *events_grpc.UpdateEventRequest) (*entities.UpdateEventRequest, error) {
	updateEventRequest := entities.UpdateEventRequest{}

	updateEventRequest.ID = storage.ID(eventGrpc.Id)
	updateEventRequest.UserID = storage.ID(eventGrpc.UserId)
	if eventGrpc.Title != nil {
		updateEventRequest.Title = eventGrpc.Title.Value
	}

	if eventGrpc.Description != nil {
		updateEventRequest.Description = eventGrpc.Description.Value
	}

	startAt, err := timestampToTime(eventGrpc.StartAt)
	updateEventRequest.StartAt = startAt
	if err != nil {
		return nil, errors.New("error converting event.startAt")
	}

	endAt, err := timestampToTime(eventGrpc.EndAt)
	updateEventRequest.EndAt = endAt
	if err != nil {
		return nil, errors.New("error converting event.endAt")
	}

	if eventGrpc.HasNotifyAt {
		if eventGrpc.NotifyAt == nil {
			updateEventRequest.NotifyAt = domain.ShouldResetTime
		} else {
			notifyAt, err := timestampToTime(eventGrpc.NotifyAt)
			updateEventRequest.NotifyAt = notifyAt
			if err != nil {
				return nil, errors.New("error converting event.notifyAt")
			}
		}
	}

	return &updateEventRequest, nil
}

func (s *Server) UpdateEvent(ctx context.Context, query *events_grpc.UpdateEventRequest) (*events_grpc.EventResponse, error) {
	updateEventRequest, err := prepareUpdateEventRequest(query)

	if err != nil {
		return nil, err
	}

	event, err := s.eventService.UpdateEvent(ctx, updateEventRequest)

	if err != nil {
		return nil, errors.Wrap(err, "problem updating event")
	}

	return createEventResponse(*event)
}

func (s *Server) DeleteEvent(ctx context.Context, query *events_grpc.DeleteEventRequest) (*events_grpc.EventResponse, error) {
	deleteEventRequest := entities.DeleteEventRequest{
		ID:     storage.ID(query.EventId),
		UserID: storage.ID(query.UserId),
	}

	event, err := s.eventService.DeleteEvent(ctx, &deleteEventRequest)

	if err != nil {
		return nil, errors.Wrap(err, "problem deleting event")
	}

	return createEventResponse(*event)
}

func (s *Server) Start(eventService domain.EventService) error {
	lsn, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return err
	}

	s.instance = grpc.NewServer()
	service := &Server{eventService: eventService}

	events_grpc.RegisterEventsServer(s.instance, service)

	fmt.Printf("Starting server on %s\n", lsn.Addr().String())
	if err := s.instance.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.instance.Stop()

	return nil
}
