package queue

import (
	"github.com/streadway/amqp"
)

type Consumer interface {
	Handle(fn func(<-chan amqp.Delivery)) error
}

type Producer interface {
	Run(msgs <-chan amqp.Publishing) error
}
