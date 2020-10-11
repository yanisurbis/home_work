package main

import (
	grpcclient "calendar/internal/client/grpc"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

/*type Sender struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	uri     string
}



func NewSender(uri string, queueName string) *Sender {
	return &Sender{
		queue: queueName,
		uri:   uri,
	}
}

func (s *Sender) connect() {
	var err error

	s.conn, err = amqp.Dial(s.uri)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer s.conn.Close()

	s.channel, err = s.conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer s.channel.Close()
}*/

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
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
	failOnError(err, "Failed to publish a message")
}
