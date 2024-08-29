package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"github.com/Nicknamezz00/gorder/pkg/broker"
	amqp "github.com/rabbitmq/amqp091-go"
)

type consumer struct {
	service OrderService
}

func NewConsumer(service OrderService) *consumer {
	return &consumer{service}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.OrderPaid, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			// Extract the headers
			//ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

			//tr := otel.Tracer("amqp")
			//_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

			o := &pb.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				_ = d.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			_, err := c.service.UpdateOrder(context.Background(), o)
			log.Printf("c.service.UpdateOrder(context.Background(), o)")
			if err != nil {
				log.Printf("failed to update order: %v", err)
				if err := broker.HandleRetry(ch, &d); err != nil {
					log.Printf("error handling retry: %v", err)
				}
				continue
			}

			//messageSpan.AddEvent("order.updated")
			//messageSpan.End()

			log.Println("order has been updated from AMQP")
			d.Ack(false)
		}
	}()

	<-forever
}
