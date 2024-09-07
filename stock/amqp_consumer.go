package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Nicknamezz00/gorder/pkg/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

type Consumer struct{}

func NewConsumer() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(
		q.Name,           // queue name
		"",               // routing key
		broker.OrderPaid, // exchange
		false,            // no-wait
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}

	go func() {
		for msg := range msgs {
			// Extract headers
			ctx := broker.ExtractAMQPHeader(context.Background(), msg.Headers)
			// Create a new span
			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

			log.Printf("Received a message: %s", msg.Body)
			orderID := string(msg.Body)

			msg.Ack(false)

			messageSpan.End()
			log.Printf("Order received: %s", orderID)
			// TODO: Do ops to stock.
		}
	}()

	log.Printf("AMQP Listening. Interrupt to exit")
	<-forever
}
