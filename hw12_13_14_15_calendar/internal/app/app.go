package app

import (
	server2 "calendar/internal/grpc/server"
	"calendar/internal/logger"
	"calendar/internal/repository"
	"calendar/internal/server"
	"context"
	"fmt"
	"log"
)

type App struct {
	repo       repository.BaseRepo
	server     server.Server
	grpcServer server2.Server
	logger     logger.Logger
}

func New(r repository.BaseRepo, s server.Server, l logger.Logger, g server2.Server) (*App, error) {
	return &App{repo: r, server: s, logger: l, grpcServer: g}, nil
}

func (a *App) Run(ctx context.Context, logPath string, dsn string) error {
	// logger
	err := a.logger.Init(logPath)
	if err != nil {
		log.Fatal(err)
	}

	// storage
	err = a.repo.Connect(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	// server
	//err = a.server.Start()
	//if err != nil {
	//	log.Fatal(err)
	//}
	err = a.grpcServer.Start(a.repo)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) error {
	fmt.Println("Shutting down...")

	if err := a.server.Stop(ctx); err != nil {
		log.Fatal(err)
	}

	if err := a.repo.Close(); err != nil {
		log.Fatal(err)
	}

	return nil
}
