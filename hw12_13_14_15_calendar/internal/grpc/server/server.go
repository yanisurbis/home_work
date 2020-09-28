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
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/golang/protobuf/ptypes/wrappers"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Server struct {
	eventService domain.EventService
	db repository.BaseRepo
}

func createEventResponse(event entities.Event) *events_grpc.Event {
	startAt, err := ptypes.TimestampProto(event.StartAt)

	if err != nil {
		log.Fatal("Type conversion error")
	}

	endAt, err := ptypes.TimestampProto(event.EndAt)

	if err != nil {
		log.Fatal("Type conversion error")
	}

	notifyAt, err := ptypes.TimestampProto(event.NotifyAt)

	if err != nil {
		// TODO: fix error handling
		log.Fatal("Type conversion error")
	}

	return &events_grpc.Event{
		Id:          uint32(event.ID),
		Title:       event.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: &wrappers.StringValue{Value: event.Description},
		UserId:      uint32(event.UserID),
		NotifyAt:    notifyAt,
	}
}

func (s *Server) GetEvents(ctx context.Context, query *events_grpc.EventsQuery, period string) (*events_grpc.EventsResponse, error) {
	from, err := ptypes.Timestamp(query.From)

	if err != nil {
		log.Fatal("Type conversion error")
	}

	getEventsRequest := entities.GetEventsRequest{
		UserID: repository.ID(query.UserId),
		Type:   period,
		From:   from,
	}

	events, err := s.eventService.GetEvents(ctx, &getEventsRequest)

	var eventsResponse []*events_grpc.Event

	for _, event := range events {
		eventsResponse = append(eventsResponse, createEventResponse(event))
	}

	return &events_grpc.EventsResponse{Events: eventsResponse}, nil
}

func (s *Server) GetEventsDay(ctx context.Context, query *events_grpc.EventsQuery) (*events_grpc.EventsResponse, error) {
	return s.GetEvents(ctx, query, domain.PeriodDay)
}

func (s *Server) GetEventsWeek(ctx context.Context, query *events_grpc.EventsQuery) (*events_grpc.EventsResponse, error) {
	return s.GetEvents(ctx, query, domain.PeriodWeek)
}

func (s *Server) GetEventsMonth(ctx context.Context, query *events_grpc.EventsQuery) (*events_grpc.EventsResponse, error) {
	return s.GetEvents(ctx, query, domain.PeriodMonth)
}

func prepareAddEventRequest(eventGrpc *events_grpc.Event) (*entities.AddEventRequest, error) {
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

	notifyAt, err := ptypes.Timestamp(eventGrpc.NotifyAt)

	if err != nil {
		return nil, errors.New("error converting event.notifyAt")
	}

	description := ""
	if eventGrpc.Description == nil {
		description = domain.DefaultEmptyString
	} else {
		description = eventGrpc.Description.Value
	}

	return &entities.AddEventRequest{
		Title:       eventGrpc.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: description,
		NotifyAt:    notifyAt,
		UserID:      repository.ID(eventGrpc.UserId),
	}, nil
}

func (s *Server) AddEvent(ctx context.Context, query *events_grpc.Event) (*empty.Empty, error) {
	addEventRequest, err := prepareAddEventRequest(query)

	if err != nil {
		return nil, err
	}

	_, err = s.eventService.AddEvent(ctx, addEventRequest)

	if err != nil {
		return nil, errors.New("problem adding event to the DB")
	}

	return &empty.Empty{}, nil
}

func prepareUpdateEventRequest(eventGrpc *events_grpc.Event) (*entities.UpdateEventRequest, error) {
	fmt.Printf("%+v\n", eventGrpc)

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

	notifyAt := domain.DefaultEmptyTime
	if eventGrpc.NotifyAt != nil {
		notifyAt, err = ptypes.Timestamp(eventGrpc.NotifyAt)
		if err != nil {
			return nil, errors.New("error converting event.notifyAt")
		}
	}

	return &entities.UpdateEventRequest{
		ID:          repository.ID(eventGrpc.Id),
		Title:       eventGrpc.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: description,
		NotifyAt:    notifyAt,
		UserID:      repository.ID(eventGrpc.UserId),
	}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, query *events_grpc.Event) (*empty.Empty, error) {
	updateEventRequest, err := prepareUpdateEventRequest(query)

	if err != nil {
		return nil, err
	}

	_, err = s.eventService.UpdateEvent(ctx, updateEventRequest)

	if err != nil {
		return nil, errors.New("problem adding event to the DB")
	}

	return &empty.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, query *events_grpc.DeleteEventRequest) (*empty.Empty, error) {
	deleteEventRequest := entities.DeleteEventRequest{
		ID:     repository.ID(query.EventId),
		UserID: repository.ID(query.UserId),
	}

	_, err := s.eventService.DeleteEvent(ctx, &deleteEventRequest)

	if err != nil {
		return nil, errors.New("problem deleting event to the DB")
	}

	return &empty.Empty{}, nil
}

func (s *Server) Start(eventService domain.EventService) error {
	lsn, err := net.Listen("tcp", "localhost:9090")

	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	service := &Server{eventService: eventService}

	events_grpc.RegisterEventsServer(server, service)

	fmt.Printf("Starting server on %s\n", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return nil
}
