package main

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/config"
	queue2 "calendar/internal/queue"
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
		var producer queue2.Producer
		// TODO: client type should be producer
		producer = rabbit.Initialize(c.Queue.ConsumerTag, "sender", c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
		_ = producer.Run(msgs)
	}()

	everyMinute := time.NewTicker(60 * time.Second)
	for {
		select {
		case <-everyMinute.C:
			notifications, err := client.GetNotifications(time.Now().Add(-2*time.Hour), time.Now())
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

			err = client.DeleteOldEvents()
			if err != nil {
				log.Println(err)
			}
		case <-ctx.Done():
			close(msgs)

			err := client.Stop()
			if err != nil {
				log.Fatal(err)
			}

			everyMinute.Stop()

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
