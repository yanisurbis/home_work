package rabbit

// docker run -d --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq:3-management
// https://github.com/rabbitmq/rabbitmq-consistent-hash-exchange
// rabbitmq-plugins enable rabbitmq_consistent_hash_exchange

import (
	"calendar/internal/queue"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"
)

type Queue struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	clientType   string
	done         chan error
	consumerTag  string
	uri          string
	exchangeName string
	exchangeType string
	queue        string
	bindingKey   string
}

const producer = "producer"
const consumer = "consumer"

func initialize(consumerTag, clientType, uri, exchangeName, exchangeType, queue, bindingKey string) *Queue {
	return &Queue{
		consumerTag:  consumerTag,
		clientType:   clientType,
		uri:          uri,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queue:        queue,
		bindingKey:   bindingKey,
		done:         make(chan error),
	}
}

func CreateProducer(consumerTag, uri, exchangeName, exchangeType, queue, bindingKey string) queue.Producer {
	return initialize(consumerTag, producer, uri, exchangeName, exchangeType, queue, bindingKey)
}

func CreateConsumer(consumerTag, uri, exchangeName, exchangeType, queue, bindingKey string) queue.Consumer {
	return initialize(consumerTag, consumer, uri, exchangeName, exchangeType, queue, bindingKey)
}

func (c *Queue) reConnect() (<-chan amqp.Delivery, error) {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2

	b := backoff.WithContext(be, context.Background())
	d := b.NextBackOff()

	for range time.After(d) {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return nil, fmt.Errorf("stop reconnecting")
		}

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

	return nil, nil
}

func (c *Queue) connect() error {
	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	c.channel, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("channel: %s", err)
	}

	go func() {
		log.Printf("closing: %s", <-c.conn.NotifyClose(make(chan *amqp.Error)))
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		c.done <- errors.New("channel closed")
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
		return fmt.Errorf("exchange Declare: %s", err)
	}

	return nil
}

// Задекларировать очередь, которую будем слушать.
func (c *Queue) announceQueue() (<-chan amqp.Delivery, error) {
	q, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return nil, fmt.Errorf("queue Declare: %s", err)
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.channel.Qos(50, 0, false)
	if err != nil {
		return nil, fmt.Errorf("error setting qos: %s", err)
	}

	// Создаём биндинг (правило маршрутизации).
	if err = c.channel.QueueBind(
		q.Name,
		c.bindingKey,
		c.exchangeName,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %s", err)
	}

	if c.clientType == consumer {
		msgs, err := c.channel.Consume(
			q.Name,
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

	return nil, nil
}

func (c *Queue) Handle(ctx context.Context, fn func(<-chan amqp.Delivery)) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("error: %v", err)
	}
	msgs, err := c.announceQueue()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for {
		go fn(msgs)

		select {
		case done := <-c.done:
			{
				if done != nil {
					msgs, err = c.reConnect()
					if err != nil {
						return fmt.Errorf("reconnecting Error: %s", err)
					}
					fmt.Println("reconnected... possibly")
				}
			}
		case <-ctx.Done():
			return c.Close()
		}
	}
}

func (c *Queue) Run(msgs <-chan amqp.Publishing) error {
	var err error
	if err = c.connect(); err != nil {
		return fmt.Errorf("error: %v", err)
	}
	_, err = c.announceQueue()
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return c.Close()
			}
			err = c.channel.Publish(
				c.exchangeName, // exchange
				c.bindingKey,   // routing key
				false,          // mandatory
				false,          // immediate
				msg,
			)
			if err != nil {
				log.Printf("failed to send a message: %v", err)
			}
		case <-c.done:
			_, err = c.reConnect()
			if err != nil {
				return fmt.Errorf("reconnecting Error: %s", err)
			}
			fmt.Println("Reconnected... possibly")
		}
	}
}

func (c *Queue) Close() error {
	return c.conn.Close()
}
