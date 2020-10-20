package main

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/config"
	"calendar/internal/queue/rabbit"
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	msgs := make(chan amqp.Publishing)
	c, _ := config.Read("./configs/local.toml")

	client := grpcclient.NewClient()
	client.Start(ctx)

	go handleSignals(cancel)
	go func() {
		producer := rabbit.CreateProducer(c.Queue.ConsumerTag, c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
		_ = producer.Run(msgs)
	}()

	everyMinute := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-everyMinute.C:
			notifications, err := client.GetNotifications(time.Now().Add(-1*time.Minute), time.Now())
			if err != nil {
				log.Println(err)
			} else {
				msg, err := json.Marshal(notifications)
				if err != nil {
					log.Println(err)
				} else {
					msgs <- amqp.Publishing{
						ContentType: "application/json",
						Body:        msg,
					}
				}
			}

			err = client.DeleteOldEvents(time.Now().Add(-1*time.Minute))
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			everyMinute.Stop()
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
	<-sigCh
}
