package main

import (
	grpcclient "calendar/internal/client/grpc"
	"calendar/internal/config"
	"calendar/internal/queue/rabbit"
	"context"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	msgs := make(chan amqp.Publishing)
	c, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	client := grpcclient.NewClient()
	err = client.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	go handleSignals(cancel)

	go func() {
		producer := rabbit.CreateProducer(c.Queue.ConsumerTag, c.Queue.URI, c.Queue.ExchangeName, c.Queue.ExchangeType, c.Queue.Queue, c.Queue.BindingKey)
		err = producer.Run(msgs)
		if err != nil {
			log.Fatal(err)
		}
	}()

	interval := time.Duration(c.Scheduler.FetchIntervalSeconds) * time.Second
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			notifications, err := client.GetNotifications(time.Now().Add(-1*interval), time.Now())
			if err != nil {
				log.Println(err)
			} else if len(notifications) > 0 {
				msg, err := json.Marshal(notifications)
				if err != nil {
					log.Println(err)
				} else {
					log.Println(time.Now().Format(time.Stamp), "sending", len(notifications), "messages")
					msgs <- amqp.Publishing{
						ContentType: "application/json",
						Body:        msg,
					}
				}
			} else if len(notifications) == 0 {
				log.Println(time.Now().Format(time.Stamp), "no notifications to send")
			}

			//err = client.DeleteOldEvents(time.Now().Add(-1 * time.Minute))
			//if err != nil {
			//	log.Println(err)
			//}
		case <-ctx.Done():
			ticker.Stop()
			close(msgs)
			err := client.Stop()
			if err != nil {
				log.Fatal(err)
			}

			return
		}
	}
}

func handleSignals(cancel context.CancelFunc) {
	defer cancel()
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh
}
