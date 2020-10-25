package app

import (
	domain "calendar/internal/domain/interfaces"
	domain2 "calendar/internal/domain/services"
	"calendar/internal/logger"
	"calendar/internal/server"
	"context"
	"fmt"
	"log"
)

type App struct {
	server     server.Server
	grpcServer server.Server
	logger     logger.Logger
	storage    domain.EventStorage
}

func New(s server.Server, grpcServer server.Server, l logger.Logger, storage domain.EventStorage) (*App, error) {
	return &App{server: s, logger: l, storage: storage, grpcServer: grpcServer}, nil
}

func (a *App) Run(ctx context.Context, logPath string, dsn string) error {
	// logger
	err := a.logger.Init(logPath)
	if err != nil {
		return err
	}

	// storage
	fmt.Println("DSN:", dsn)
	err = a.storage.Connect(ctx, dsn)
	if err != nil {
		return err
	}

	// service
	eventService := domain2.EventService{
		EventStorage: a.storage,
	}

	// http server
	go func() {
		err = a.server.Start(eventService)
		if err != nil {
			// TODO: handle fatal
			log.Fatal(err)
		}
	}()

	// grpc server
	err = a.grpcServer.Start(eventService)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Println("Shutting down...")

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
