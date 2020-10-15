package main

import (
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

	consumer = rabbit.Initialize("events_consumer", "consumer", "amqp://guest:guest@localhost:5672/", "exchange", "fanout", "events_notifications", "events")
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
