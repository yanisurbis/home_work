package main

import (
	"calendar/internal/config"
	"calendar/internal/domain/entities"
	"calendar/internal/queue/rabbit"
	"context"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	c, _ := config.Read("./configs/local.toml")

	consumer := rabbit.CreateConsumer(c.Queue.ConsumerTag, c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
	go handleSignals(cancel)
	_ = consumer.Handle(ctx, func(msgs <-chan amqp.Delivery) {
		for {
			select {
			case msg, ok := <-msgs:
				if ok == false {
					return
				}
				var notifications []entities.Notification
				err := json.Unmarshal(msg.Body, &notifications)
				if err != nil {
					failOnError(err, "Error unmarshalling json")
				}

				fmt.Println("%+v", notifications)
			}
		}
	})

}

func handleSignals(cancel context.CancelFunc) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
}
