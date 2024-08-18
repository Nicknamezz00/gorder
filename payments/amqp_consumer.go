package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{
		service: service,
	}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	var forever chan struct{}
	go func() {
		for msg := range msgs {
			log.Printf("received message: %s", msg.Body)
			o := &pb.Order{}
			if err := json.Unmarshal(msg.Body, o); err != nil {
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}
			link, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("failed to create payment link: %v", err)
				continue
			}
			log.Printf("payment link: %s", link)
		}
	}()
	<-forever
}
