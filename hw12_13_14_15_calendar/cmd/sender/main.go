package main

import (
	"calendar/internal/config"
	"calendar/internal/domain/entities"
	"calendar/internal/queue/rabbit"
	"calendar/internal/storage/sql"
	"context"
	"encoding/json"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10*time.Second)
	storage := new(sql.Repo)
	err = storage.Connect(context.Background(), c.PSQL.DSN)
	if err != nil {
		log.Fatal(err)
	}

	consumer := rabbit.CreateConsumer(
		c.Queue.ConsumerTag,
		c.Queue.URI,
		c.Queue.ExchangeName,
		c.Queue.ExchangeType,
		c.Queue.Queue,
		c.Queue.BindingKey,
	)
	go handleSignals(cancel)

	err = consumer.Handle(ctx, func(msgs <-chan amqp.Delivery) {
		for msg := range msgs {
			var notifications []entities.Notification
			err := json.Unmarshal(msg.Body, &notifications)

			if err != nil {
				log.Println(err)
				continue
			}

			log.Println("received", len(notifications), "notifications")

			//if os.Getenv("ENV") == "TEST" {
			//	err = storage.AddNotifications(notifications)
			//	if err != nil {
			//		log.Println(err)
			//	}
			//}
			err = storage.AddNotifications(notifications)
			if err != nil {
				log.Println(err)
			}
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func handleSignals(cancel context.CancelFunc) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	signal.Stop(sigCh)
	<-sigCh
}
