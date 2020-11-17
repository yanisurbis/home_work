package app

import (
	"calendar/internal/config"
	domain "calendar/internal/domain/interfaces"
	domain2 "calendar/internal/domain/services"
	"calendar/internal/logger"
	"calendar/internal/server"
	"context"
	"log"
	"time"
)

type App struct {
	server     server.Server
	grpcServer server.Server
	logger     logger.Logger
	storage    domain.EventStorage
}

func New(
	s server.Server,
	grpcServer server.Server,
	l logger.Logger,
	storage domain.EventStorage,
) (*App, error) {
	return &App{server: s, logger: l, storage: storage, grpcServer: grpcServer}, nil
}

func (a *App) Run(ctx context.Context, config *config.Config) error {
	// logger
	err := a.logger.Init(config.Logger.Path)
	if err != nil {
		return err
	}

	// storage
	err = a.storage.Connect(ctx, config.PSQL.DSN)
	if err != nil {
		return err
	}

	events, err := a.storage.GetEventsMonth(1, time.Now().Add(-1 * time.Hour))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("events: %v", events)

	// service
	eventService := domain2.EventService{
		EventStorage: a.storage,
	}

	// http server
	go func() {
		err = a.server.Start(eventService, config.HTTPServer.Address)
		if err != nil {
			// TODO: handle fatal
			log.Fatal(err)
		}
	}()

	// grpc server
	err = a.grpcServer.Start(eventService, config.GRPCServer.Address)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	log.Println("Shutting down...")

	if err := a.server.Stop(ctx); err != nil {
		return err
	}

	if err := a.grpcServer.Stop(ctx); err != nil {
		return err
	}

	if err := a.storage.Close(); err != nil {
		return err
	}

	return nil
}
