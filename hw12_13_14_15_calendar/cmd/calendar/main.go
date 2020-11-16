package main

import (
	"calendar/internal/app"
	"calendar/internal/config"
	"calendar/internal/logger"
	grpcserver "calendar/internal/server/grpc/server"
	httpserver "calendar/internal/server/http"
	"calendar/internal/storage/sql"
	"context"
	"log"
	"os"
	"os/signal"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	s := new(httpserver.Instance)
	grpcServer := new(grpcserver.Server)
	l := new(logger.Instance)
	storage := new(sql.Repo)

	a, err := app.New(s, grpcServer, l, storage)
	if err != nil {
		log.Fatal(err)
	}

	go handleSignals(ctx, cancel, a)

	if err := a.Run(ctx, conf); err != nil {
		log.Fatal(err)
	}
}

func handleSignals(ctx context.Context, cancel context.CancelFunc, app *app.App) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Stop(sigCh)
	<-sigCh
	err := app.Stop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
