package main

import (
	grpcclient "calendar/internal/client/grpc"
	queue2 "calendar/internal/queue"
	"calendar/internal/queue/rabbit"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	var producer queue2.Producer

	producer = rabbit.Initialize("events_consumer", "sender", "amqp://guest:guest@localhost:5672/", "exchange", "fanout", "events_notifications", "events")

	channel := make(chan amqp.Publishing)

	go func() {
		_ = producer.Run(channel)
	}()

	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			notifications := grpcclient.GetNotifications()

			msg, err := json.Marshal(notifications)
			failOnError(err, "Couldn't serialize notifications")
			channel <- amqp.Publishing{
				ContentType: "application/json",
				Body:        msg,
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
