package main

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/queue/rabbit"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	queue := rabbit.NewQueue("events_consumer", "sender", "amqp://guest:guest@localhost:5672/", "exchange", "fanout", "events_notifications", "events")

	channel := make(chan amqp.Publishing)

	go queue.Run(channel)
	notifications := grpcclient.GetNotifications()

	msg, err := json.Marshal(notifications)
	failOnError(err, "Couldn't serialize notifications")
	channel <- amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	}

	emptyChannel := make(chan int)

	<-emptyChannel

	/*failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"notifications", // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)

	notifications := grpcclient.GetNotifications()

	failOnError(err, "Failed to declare a queue")
	msg, err := json.Marshal(notifications)
	failOnError(err, "Couldn't serialize notifications")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	log.Printf(" [x] Sent %s", msg)
	failOnError(err, "Failed to publish a message")*/
}
