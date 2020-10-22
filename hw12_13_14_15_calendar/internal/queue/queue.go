package queue

import (
	"context"

	"github.com/streadway/amqp"
)

type Consumer interface {
	Handle(ctx context.Context, fn func(<-chan amqp.Delivery)) error
	Close() error
}

type Producer interface {
	Run(msgs <-chan amqp.Publishing) error
	Close() error
}
