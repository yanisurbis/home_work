package main

import (
	"calendar/internal/app"
	"calendar/internal/config"
	"calendar/internal/logger"
	"calendar/internal/protobufs"
	"calendar/internal/repository/postgres"
	"calendar/internal/server"
	"context"
	"flag"
	"log"

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

	//c, _ := config.Read(args.configPath)
	c, _ := config.Read("./configs/local.toml")

	r := new(postgres.Repo)
	s := new(server.Instance)
	gs := new(protobufs.Server)
	l := new(logger.Instance)

	// TODO why deref gs
	a, err := app.New(r, s, l, *gs)
	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(ctx, c.Logger.Path, c.PSQL.DSN); err != nil {
		log.Fatal(err)
	}
}
