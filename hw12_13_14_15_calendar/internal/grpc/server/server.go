//TODO: Rename to server_grpc
package server

import (
	"calendar/internal/grpc/events_grpc"
	"calendar/internal/repository"
	"context"
	"errors"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type Server struct {
	db repository.BaseRepo
}

func createEventResponse(event repository.Event) *events_grpc.Event {
	start_at, err := ptypes.TimestampProto(event.StartAt)

	if err != nil {
		log.Fatal("Type conversion error")
	}

	end_at, err := ptypes.TimestampProto(event.EndAt)

	if err != nil {
		log.Fatal("Type conversion error")
	}

	notify_at, err := ptypes.TimestampProto(event.NotifyAt)

	if err != nil {
		// TODO: fix error handling
		log.Fatal("Type conversion error")
	}

	return &events_grpc.Event{
		Id:          uint32(event.ID),
		Title:       event.Title,
		StartAt:     start_at,
		EndAt:       end_at,
		Description: event.Description,
		UserId:      uint32(event.UserID),
		NotifyAt:    notify_at,
	}
}

func getEvents(query *events_grpc.EventsQuery, cb func(userID repository.ID, from time.Time) ([]repository.Event, error)) (*events_grpc.EventsResponse, error) {
	t, err := ptypes.Timestamp(query.From)

	if err != nil {
		log.Fatal("Type conversion error")
	}

	events, err := cb(repository.ID(query.UserId), t)

	var eventsResponse []*events_grpc.Event

	for _, event := range events {
		eventsResponse = append(eventsResponse, createEventResponse(event))
	}

	return &events_grpc.EventsResponse{Events: eventsResponse}, nil
}

func (s *Server) GetEventsDay(ctx context.Context, query *events_grpc.EventsQuery) (*events_grpc.EventsResponse, error) {
	return getEvents(query, s.db.GetEventsDay)
}

func (s *Server) GetEventsWeek(ctx context.Context, query *events_grpc.EventsQuery) (*events_grpc.EventsResponse, error) {
	return getEvents(query, s.db.GetEventsWeek)
}

func (s *Server) GetEventsMonth(ctx context.Context, query *events_grpc.EventsQuery) (*events_grpc.EventsResponse, error) {
	return getEvents(query, s.db.GetEventsMonth)
}

func convertEvent(eventGrpc *events_grpc.Event) (*repository.Event, error) {
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

	return &repository.Event{
		ID:          repository.ID(eventGrpc.Id),
		Title:       eventGrpc.Title,
		StartAt:     startAt,
		EndAt:       endAt,
		Description: eventGrpc.Description,
		UserID:      repository.ID(eventGrpc.UserId),
		NotifyAt:    notifyAt,
	}, nil
}

func (s *Server) AddEvent(ctx context.Context, query *events_grpc.Event) (*empty.Empty, error) {
	event, err := convertEvent(query)

	if err != nil {
		// TODO: check nil handling
		return nil, err
	}

	err = s.db.AddEvent(*event)

	if err != nil {
		return nil, errors.New("problem adding event to the DB")
	}

	return &empty.Empty{}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, query *events_grpc.Event) (*empty.Empty, error) {
	event, err := convertEvent(query)

	if err != nil {
		return nil, err
	}

	err = s.db.UpdateEvent(*event)

	if err != nil {
		return nil, errors.New("problem updating event to the DB")
	}

	return &empty.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, query *events_grpc.DeleteEventRequest) (*empty.Empty, error) {
	err := s.db.DeleteEvent(repository.ID(query.UserId), repository.ID(query.EventId))

	if err != nil {
		return nil, errors.New("problem deleting event to the DB")
	}

	return &empty.Empty{}, nil
}

func (s *Server) Start(r repository.BaseRepo) error {
	lsn, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	service := &Server{db: r}

	events_grpc.RegisterEventsServer(server, service)

	fmt.Printf("Starting server on %s\n", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}
