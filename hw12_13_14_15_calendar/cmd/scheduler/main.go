package main

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/config"
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
	//grpcclient.DeleteOldEvents()
	channel := make(chan amqp.Publishing)
	c, _ := config.Read("./configs/local.toml")
	client := grpcclient.NewClient()
	client.Start()

	go func() {
		var producer queue2.Producer
		// TODO: client type should be producer
		producer = rabbit.Initialize(c.Queue.ConsumerTag, "sender", c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
		_ = producer.Run(channel)
	}()

	everyMinute := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-everyMinute.C:
			notifications := client.GetNotifications(time.Now().Add(-2*time.Hour), time.Now())

			msg, err := json.Marshal(notifications)
			failOnError(err, "Couldn't serialize notifications")
			channel <- amqp.Publishing{
				ContentType: "application/json",
				Body:        msg,
			}

			client.DeleteOldEvents()
		case <-quit:
			everyMinute.Stop()
			return
		}
	}
}
