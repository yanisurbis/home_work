package main

import (
	"calendar/internal/domain/entities"
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
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var notifications []entities.Notification

			err := json.Unmarshal(d.Body, &notifications)
			if err != nil {
				log.Fatal("Error unmarshalling json")
			}

			fmt.Println("%+v", notifications)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
