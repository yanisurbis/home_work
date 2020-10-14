package consumer

// docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management
// https://github.com/rabbitmq/rabbitmq-consistent-hash-exchange
// rabbitmq-plugins enable rabbitmq_consistent_hash_exchange

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"
	"github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"
)

// Consumer ...
type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	done         chan error
	consumerTag  string
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
	maxInterval  time.Duration
}

const defaultMaxInt = time.Second * 15

type Opt func(consumer *Consumer)

func WithMaxInterval(interval time.Duration) Opt {
	return func(consumer *Consumer) {
		consumer.maxInterval = interval
	}
}

func NewConsumer(consumerTag, uri, exchangeName, exchangeType, queue, bindingKey string, opts ...Opt) *Consumer {
	c := &Consumer{
		consumerTag:  consumerTag,
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queue:        queue,
		bindingKey:   bindingKey,
		done:         make(chan error),
		maxInterval:  defaultMaxInt,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Consumer) reConnect() (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

		select {
		case <-time.After(d):
			if err := c.connect(); err != nil {
				log.Printf("could not connect in reconnect call: %+v", err)
				continue
			}
			msgs, err := c.announceQueue()
			if err != nil {
				fmt.Printf("Couldn't connect: %+v", err)
				continue
			}

			return msgs, nil
		}
	}
}

func (c *Consumer) connect() error {

	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		c.done <- errors.New("Channel Closed")
	}()

	if err = c.channel.ExchangeDeclare(
		c.exchangeName,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("Exchange Declare: %s", err)
	}

	return nil
}

// Задекларировать очередь, которую будем слушать.
func (c *Consumer) announceQueue() (<-chan amqp.Delivery, error) {
	queue, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("Queue Declare: %s", err)
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("Error setting qos: %s", err)
	}

	// Создаём биндинг (правило маршрутизации).
	if err = c.channel.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("Queue Bind: %s", err)
	}

	msgs, err := c.channel.Consume(
		queue.Name,
		c.consumerTag,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Queue Consume: %s", err)
	}

	return msgs, nil
}

func (c *Consumer) Handle(fn func(<-chan amqp.Delivery), threads int) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	msgs, err := c.announceQueue()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	for {
		for i := 0; i < threads; i++ {
			go fn(msgs)
		}

		if <-c.done != nil {
			msgs, err = c.reConnect()
			if err != nil {
				return fmt.Errorf("Reconnecting Error: %s", err)
			}
		}
		fmt.Println("Reconnected... possibly")
	}
}

func (c *Consumer) Run(mesgs <-chan amqp.Publishing) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("Error: %v", err)
	}
	_, err = c.announceQueue()
	if err != nil {
		return fmt.Errorf("Error: %v", err)
	}

	for {
		select {
		case msg := <-mesgs:
			err = c.channel.Publish(
				"",     // exchange
				"TODO", // routing key
				false,  // mandatory
				false,  // immediate
				msg,
			)
		case <-c.done:
			_, err = c.reConnect()
			if err != nil {
				return fmt.Errorf("Reconnecting Error: %s", err)
			}
		}
		fmt.Println("Reconnected... possibly")
	}
}