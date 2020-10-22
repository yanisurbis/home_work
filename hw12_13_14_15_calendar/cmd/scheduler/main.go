package main

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/config"
	"calendar/internal/queue/rabbit"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/streadway/amqp"
)

const interval = 5 * time.Second

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	msgs := make(chan amqp.Publishing)
	c, _ := config.Read("./configs/local.toml")

	client := grpcclient.NewClient()
	err := client.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go handleSignals(cancel)
	go func() {
		producer := rabbit.CreateProducer(c.Queue.ConsumerTag, c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
		_ = producer.Run(msgs)
	}()

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			notifications, err := client.GetNotifications(time.Now().Add(-1*interval), time.Now())
			//notifications, err := client.GetNotifications(time.Now().Add(-2*time.Hour), time.Now())
			if err != nil {
				log.Println(err)
			} else {
				msg, err := json.Marshal(notifications)
				if err != nil {
					log.Println(err)
				} else {
					fmt.Println(time.Now().Format(time.Stamp), "sending", len(notifications), "messages")
					msgs <- amqp.Publishing{
						ContentType: "application/json",
						Body:        msg,
					}
				}
			}

			//err = client.DeleteOldEvents(time.Now().Add(-1*time.Minute))
			//if err != nil {
			//	log.Println(err)
			//}
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
	<-sigCh
}
