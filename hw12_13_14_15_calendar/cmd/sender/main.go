package main

import (
	"calendar/internal/config"
	"calendar/internal/domain/entities"
	queue2 "calendar/internal/queue"
	"calendar/internal/queue/rabbit"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	var consumer queue2.Consumer
	c, _ := config.Read("./configs/local.toml")

	consumer = rabbit.Initialize(c.Queue.ConsumerTag, "consumer", c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
	_ = consumer.Handle(func(msgs <-chan amqp.Delivery) {
		for {
			select {
			case msg := <-msgs:
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
