package main

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/config"
	"calendar/internal/queue/rabbit"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go handleSignals(cancel)

	msgs := make(chan amqp.Publishing)
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := grpcclient.NewClient()
	err = client.Start(ctx, c.GRPCServer)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		producer := rabbit.CreateProducer(
			c.Queue.ConsumerTag,
			c.Queue.URI,
			c.Queue.ExchangeName,
			c.Queue.ExchangeType,
			c.Queue.Queue,
			c.Queue.BindingKey,
		)
		err = producer.Run(msgs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	interval := time.Duration(c.Scheduler.FetchIntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			notifications, err := client.GetNotifications(time.Now().Add(-1*interval), time.Now())

			if err != nil {
				log.Println(err)
				continue
			}

			if len(notifications) > 0 {
				msg, err := json.Marshal(notifications)
				if err != nil {
					log.Println(err)
					continue
				}

				log.Println(
					time.Now().Format(time.Stamp),
					"sending",
					len(notifications),
					"messages",
				)
				select {
				case <-ctx.Done():
					continue
				case msgs <- amqp.Publishing{
					ContentType: "application/json",
					Body:        msg,
				}:
				}
			}

			// TODO: create new table to track events which were sent
			yearAgo := time.Now().Add(-24 * 30 * 12 * time.Hour)
			err = client.DeleteOldEvents(yearAgo)
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			ticker.Stop()
			close(msgs)
			err := client.Stop()
			if err != nil {
				log.Fatal(err)
			}

			return
		}
	}
}

func handleSignals(cancel context.CancelFunc) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Stop(sigCh)
	<-sigCh
}
