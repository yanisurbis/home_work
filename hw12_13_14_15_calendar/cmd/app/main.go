package main

import (
	"calendar/internal/app"
	"calendar/internal/config"
	grpcserver "calendar/internal/grpc/server"
	"calendar/internal/logger"
	"calendar/internal/repository/postgres"
	httpserver "calendar/internal/server/http"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Args struct {
	configPath string
}

func getArgs() *Args {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	args := Args{
		configPath: *configPath,
	}

	return &args
}

func main() {
	//args := getArgs()

	ctx, cancel := context.WithCancel(context.Background())

	c, _ := config.Read("./configs/local.toml")
	//c, _ := config.Read(args.configPath)

	s := new(httpserver.Instance)
	grpcServer := new(grpcserver.Server)
	l := new(logger.Instance)
	storage := new(postgres.Repo)

	a, err := app.New(s, grpcServer, l, storage)

	if err != nil {
		log.Fatal(err)
	}

	go handleSignals(ctx, cancel, a)

	if err := a.Run(ctx, c.Logger.Path, c.PSQL.DSN); err != nil {
		log.Fatal(err)
	}
}

func handleSignals(ctx context.Context, cancel context.CancelFunc, app *app.App) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	err := app.Stop(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
