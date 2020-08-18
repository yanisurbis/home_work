package server

import (
	"calendar/internal/grpc/events_grpc"
	"calendar/internal/repository"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"log"
	"net"
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

func (s *Server) GetEventsDay(ctx context.Context, query *events_grpc.EventsQuery) (*events_grpc.EventsResponse, error) {
	fmt.Println("Hello from grpc server")

	time, err := ptypes.Timestamp(query.From)

	if err != nil {
		log.Fatal("Type conversion error")
	}

	events, err := s.db.GetEventsDay(repository.ID(query.UserId), time)

	var eventsResponse []*events_grpc.Event

	for _, event := range events {
		eventsResponse = append(eventsResponse, createEventResponse(event))
	}

	return &events_grpc.EventsResponse{Events: eventsResponse}, nil
}

func (s *Server) Start(r repository.BaseRepo) error {
	lsn, err := net.Listen("tcp", "localhost:8080")

	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	service := &Server{db: r}

	events_grpc.RegisterEventsServer(server, service)

	fmt.Println("Starting server on %s", lsn.Addr().String())
	if err := server.Serve(lsn); err != nil {
		log.Fatal(err)

		return err
	}

	return nil
}
