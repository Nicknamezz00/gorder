package broker

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MaxRetryCount = 3
	DLQ           = "dlq_main"
	DLX           = "dlx_main"
)

const (
	exchangeTypeDirect = "direct"
	exchangeTypeFanout = "fanout"
)

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	addr := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)
	conn, err := amqp.Dial(addr)
	if err != nil {
		log.Fatal(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare(OrderCreated, exchangeTypeDirect, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = ch.ExchangeDeclare(OrderPaid, exchangeTypeFanout, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = createDLQAndDLX(ch)
	if err != nil {
		log.Fatal(err)
	}
	return ch, conn.Close
}

func createDLQAndDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare("main_queue", true, false, false, false, nil)
	if err != nil {
		return err
	}

	// DLX
	if err := ch.ExchangeDeclare(DLX, exchangeTypeFanout, true, false, false, false, nil); err != nil {
		return err
	}
	// Bind main queue to DLX
	if err = ch.QueueBind(q.Name, "", DLX, false, nil); err != nil {
		return err
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(DLQ, true, false, false, false, nil)
	if err != nil {
		return err
	}

	return err
}

func HandleRetry(ch *amqp.Channel, d *amqp.Delivery) error {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}

	retryCount, ok := d.Headers["x-retry-count"].(int64)
	if !ok {
		retryCount = 0
	}
	retryCount++
	d.Headers["x-retry-count"] = retryCount

	log.Printf("Retrying message %s, retry count: %d", d.Body, retryCount)

	if retryCount >= MaxRetryCount {
		log.Printf("Moving message to DLQ %s", DLQ)

		return ch.PublishWithContext(context.Background(), "", DLQ, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Headers:      d.Headers,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		})
	}

	time.Sleep(time.Second * time.Duration(retryCount))

	return ch.PublishWithContext(context.Background(), d.Exchange, d.RoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Headers:      d.Headers,
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	},
	)
}
