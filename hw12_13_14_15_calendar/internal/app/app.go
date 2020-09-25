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
	server server.Server
	logger logger.Logger
	storage domain.EventStorage
}

func New(s server.Server, l logger.Logger, storage domain.EventStorage) (*App, error) {
	return &App{server: s, logger: l, storage: storage}, nil
}

func (a *App) Run(ctx context.Context, logPath string, dsn string) error {
	// logger
	err := a.logger.Init(logPath)
	if err != nil {
		log.Fatal(err)
	}

	// storage
	err = a.storage.Connect(ctx, dsn)
	log.Println(err)
	if err != nil {
		log.Fatal(err)
	}

	// service
	eventService := domain2.EventService{
		EventStorage: a.storage,
	}

	err = a.server.Start(eventService)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Println("Shutting down...")

	if err := a.server.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	if err := a.storage.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}
