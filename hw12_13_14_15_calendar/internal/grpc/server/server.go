//TODO: Rename to server_grpc
package server

import (
	"calendar/internal/domain/entities"
	domain "calendar/internal/domain/services"
	"calendar/internal/grpc/events_grpc"
	"calendar/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"net"
	"time"
)

type Server struct {
	eventService domain.EventService
	db           repository.BaseRepo
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
	from, err := ptypes.Timestamp(query.From)
	if err != nil {
		return nil, errors.New("from conversion error")
	}

	getEventsRequest := entities.GetEventsRequest{
		UserID: repository.ID(query.UserId),
		Type:   period,
		From:   from,
	}

	events, err := s.eventService.GetEvents(ctx, &getEventsRequest)

	var eventsResponse []*events_grpc.EventResponse

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
	// TODO: check how to handle errors
	// TODO: check error handling with real errors
	// TODO: memory error if we stop the server

	startAt, err := ptypes.Timestamp(eventGrpc.StartAt)

	if err != nil {
		return nil, errors.New("error converting event.startAt")
	}

	endAt, err := ptypes.Timestamp(eventGrpc.EndAt)

	if err != nil {
		return nil, errors.New("error converting event.endAt")
	}

	notifyAt := time.Time{}
	if eventGrpc.NotifyAt != nil {
		notifyAt, err = ptypes.Timestamp(eventGrpc.NotifyAt)
		if err != nil {
			return nil, errors.New("error converting event.notifyAt")
		}
	}

	return &entities.AddEventRequest{
		Title:       eventGrpc.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: eventGrpc.Description,
		NotifyAt:    notifyAt,
		UserID:      repository.ID(eventGrpc.UserId),
	}, nil
}

func (s *Server) AddEvent(ctx context.Context, query *events_grpc.AddEventRequest) (*events_grpc.EventResponse, error) {
	addEventRequest, err := prepareAddEventRequest(query)

	if err != nil {
		return nil, err
	}

	event, err := s.eventService.AddEvent(ctx, addEventRequest)

	if err != nil {
		return nil, errors.New("problem adding event to the DB")
	}

	return createEventResponse(*event)
}

func prepareUpdateEventRequest(eventGrpc *events_grpc.UpdateEventRequest) (*entities.UpdateEventRequest, error) {
	// TODO: wrap start_at and end_at in if
	fmt.Printf("%+v\n", eventGrpc)

	title := domain.DefaultEmptyString
	if eventGrpc.Description != nil {
		title = eventGrpc.Title.Value
	}

	startAt, err := ptypes.Timestamp(eventGrpc.StartAt)
	if err != nil {
		return nil, errors.New("error converting event.startAt")
	}

	endAt, err := ptypes.Timestamp(eventGrpc.EndAt)
	if err != nil {
		return nil, errors.New("error converting event.endAt")
	}

	description := domain.DefaultEmptyString
	if eventGrpc.Description != nil {
		description = eventGrpc.Description.Value
	}

	notifyAt := time.Time{}
	if eventGrpc.HasNotifyAt {
		if eventGrpc.NotifyAt == nil {
			notifyAt = domain.DefaultEmptyTime
		} else {
			notifyAt, err = ptypes.Timestamp(eventGrpc.NotifyAt)
			if err != nil {
				return nil, errors.New("error converting event.notifyAt")
			}
		}
	}

	return &entities.UpdateEventRequest{
		ID:          repository.ID(eventGrpc.Id),
		Title:       title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: description,
		NotifyAt:    notifyAt,
		UserID:      repository.ID(eventGrpc.UserId),
	}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, query *events_grpc.UpdateEventRequest) (*events_grpc.EventResponse, error) {
	updateEventRequest, err := prepareUpdateEventRequest(query)

	if err != nil {
		return nil, err
	}

	event, err := s.eventService.UpdateEvent(ctx, updateEventRequest)

	if err != nil {
		return nil, errors.New("problem adding event to the DB")
	}

	return createEventResponse(*event)
}

func (s *Server) DeleteEvent(ctx context.Context, query *events_grpc.DeleteEventRequest) (*events_grpc.EventResponse, error) {
	deleteEventRequest := entities.DeleteEventRequest{
		ID:     repository.ID(query.EventId),
		UserID: repository.ID(query.UserId),
	}

	event, err := s.eventService.DeleteEvent(ctx, &deleteEventRequest)

	if err != nil {
		return nil, errors.New("problem deleting event to the DB")
	}

	return createEventResponse(*event)
}

func (s *Server) Start(eventService domain.EventService) error {
	lsn, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	service := &Server{eventService: eventService}

	events_grpc.RegisterEventsServer(server, service)

	fmt.Printf("Starting server on %s\n", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}
