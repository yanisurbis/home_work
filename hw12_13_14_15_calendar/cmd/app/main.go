package main

import (
	"calendar/internal/app"
	"calendar/internal/config"
	"calendar/internal/logger"
	"calendar/internal/repository/postgres"
	"calendar/internal/server"
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
	defer cancel()

	c, _ := config.Read("./configs/local.toml")
	//c, _ := config.Read(args.configPath)

	r := new(postgres.Repo)
	s := new(server.Instance)
	l := new(logger.Instance)

	a, err := app.New(r, s, l)

	if err != nil {
		log.Fatal(err)
	}

	go handleSignals(a, cancel)

	if err := a.Run(ctx, c.Logger.Path, c.PSQL.DSN); err != nil {
		log.Fatal(err)
	}
}

func handleSignals(app *app.App, cancel context.CancelFunc) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	err := app.Stop()
	if err != nil {
		log.Fatal(err)
	}
}
