package main

import (
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
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello"
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")
}
