package broker

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel"
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

type AmqpHeaderCarrier map[string]interface{}

func (a AmqpHeaderCarrier) Get(k string) string {
	value, ok := a[k]
	if !ok {
		return ""
	}

	return value.(string)
}

func (a AmqpHeaderCarrier) Set(k string, v string) {
	a[k] = v
}

func (a AmqpHeaderCarrier) Keys() []string {
	keys := make([]string, len(a))
	i := 0

	for k := range a {
		keys[i] = k
		i++
	}

	return keys
}

func InjectAMQPHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(AmqpHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func ExtractAMQPHeader(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, AmqpHeaderCarrier(headers))
}
