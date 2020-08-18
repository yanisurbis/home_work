package app

import (
	"calendar/internal/logger"
	"calendar/internal/protobufs"
	"calendar/internal/repository"
	"calendar/internal/server"
	"context"
	"log"
)

type App struct {
	repo       repository.BaseRepo
	server     server.Server
	grpcServer protobufs.Server
	logger     logger.Logger
}

func New(r repository.BaseRepo, s server.Server, l logger.Logger, g protobufs.Server) (*App, error) {
	return &App{repo: r, server: s, logger: l, grpcServer: g}, nil
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
	/*	err = a.repo.Connect(ctx, dsn)
		if err != nil {
			log.Fatal(err)
		}*/

	err = a.grpcServer.Start(a.repo)

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
