package main

import (
	"calendar/internal/app"
	"calendar/internal/config"
	server2 "calendar/internal/grpc/server"
	"calendar/internal/logger"
	"calendar/internal/repository/postgres"
	"calendar/internal/server"
	"context"
	"flag"
	"fmt"
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
	gs := new(server2.Server)
	l := new(logger.Instance)

	// TODO why deref gs
	fmt.Println("Hello")
	a, err := app.New(r, s, l, *gs)
	if err != nil {
		log.Fatal(err)
	}

	if err := a.Run(ctx, c.Logger.Path, c.PSQL.DSN); err != nil {
		log.Fatal(err)
	}
}
