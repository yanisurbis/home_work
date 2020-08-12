package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Args struct {
	host    string
	port    string
	timeout int
}

func main() {
	args := getArgs()
	ctx, cancel := context.WithCancel(context.Background())

	tc := NewTelnetClient(args.host+":"+args.port, time.Duration(args.timeout)*time.Second, os.Stdin, os.Stdout)

	if err := tc.Connect(); err != nil {
		log.Fatal(err)
	}

	go func(ctx context.Context, tc TelnetClient, cancel context.CancelFunc) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := tc.Receive()
				if err != nil {
					cancel()
					return
				}
			}
		}
	}(ctx, tc, cancel)

	go func(ctx context.Context, tc TelnetClient, cancel context.CancelFunc) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				err := tc.Send()
				if err != nil {
					cancel()
					return
				}
			}
		}
	}(ctx, tc, cancel)

	go handleSignals(tc, cancel)

	<-ctx.Done()
}

func handleSignals(tc io.Closer, cancel context.CancelFunc) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
	err := tc.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func getArgs() *Args {
	timeout := flag.String("timeout", "10s", "connection timeout")
	flag.Parse()
	otherArgs := flag.Args()

	if len(otherArgs) != 2 {
		log.Fatal("Please specify both host and port")
	}

	args := Args{
		host: otherArgs[0],
		port: otherArgs[1],
		timeout: func() int {
			timeoutStr := *timeout
			timeoutInt, err := strconv.Atoi(timeoutStr[:len(timeoutStr)-1])
			if err != nil {
				log.Fatal("Timeout format is not correct")
			}
			return timeoutInt
		}(),
	}

	if args.host == "" {
		log.Fatal("Specify 1st parameter: the host")
	}
	if args.port == "" {
		log.Fatal("Specify 2nd parameter: the port")
	}

	return &args
}
