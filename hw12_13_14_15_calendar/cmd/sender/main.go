package main

import (
	"calendar/internal/config"
	"calendar/internal/domain/entities"
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

	c, err := config.Read("./configs/local.toml")
	if err != nil {
		log.Fatal(err)
	}

	consumer := rabbit.CreateConsumer(c.Queue.ConsumerTag, c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
	go handleSignals(cancel)

	err = consumer.Handle(ctx, func(msgs <-chan amqp.Delivery) {
		for msg := range msgs {
			var notifications []entities.Notification
			err := json.Unmarshal(msg.Body, &notifications)
			if err != nil {
				log.Println(err)
			//	TODO: add more channels?
			} else {
				if len(notifications) != 0 {
					for _, notification := range notifications {
						log.Println(time.Now().Format(time.Stamp), notification.EventID, notification.EventTitle, notification.StartAt)
					}
				} else {
					log.Println(time.Now().Format(time.Stamp), "zero events received")
				}
			}
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
}
