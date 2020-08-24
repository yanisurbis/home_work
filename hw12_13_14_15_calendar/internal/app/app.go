package app

import (
	"calendar/internal/logger"
	"calendar/internal/repository"
	"calendar/internal/server"
	"context"
	"fmt"
	"log"
)

type App struct {
	repo   repository.BaseRepo
	server server.Server
	logger logger.Logger
}

func New(r repository.BaseRepo, s server.Server, l logger.Logger) (*App, error) {
	return &App{repo: r, server: s, logger: l}, nil
}

func (a *App) Run(ctx context.Context, logPath string, dsn string) error {
	// logger
	err := a.logger.Init(logPath)
	if err != nil {
		log.Fatal(err)
	}

	// server
	err = a.server.Start()
	if err != nil {
		log.Fatal(err)
	}

	// storage
	err = a.repo.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *App) Stop() error {
	fmt.Println("Stopped")
	return nil
}