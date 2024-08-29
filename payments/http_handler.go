package main

import (
	"context"
	"encoding/json"
	"fmt"
	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
	"go.opentelemetry.io/otel"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type PaymentHandler struct {
	channel *amqp.Channel
}

func NewPaymentHandler(channel *amqp.Channel) *PaymentHandler {
	return &PaymentHandler{channel: channel}
}

func (h *PaymentHandler) registerRoutes(router *http.ServeMux) {
	router.HandleFunc("/webhook", h.handleWebhook)
}

func (h *PaymentHandler) handleWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), endpointStripeSecret)
	log.Printf("got event: %v", event)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	switch event.Type {
	case stripe.EventTypeCheckoutSessionCompleted:
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Fan out, broadcast the paid event.
		if session.PaymentStatus == stripe.CheckoutSessionPaymentStatusPaid {
			log.Printf("payment for checkout session %v succeeded!", session.ID)

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			marshalledOrder, err := json.Marshal(&pb.Order{
				ID:          session.Metadata["orderID"],
				CustomerID:  session.Metadata["customerID"],
				Status:      string(stripe.CheckoutSessionPaymentStatusPaid),
				PaymentLink: "",
			})
			if err != nil {
				log.Fatal(err.Error())
			}

			tr := otel.Tracer("amqp")
			amqpContext, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - publish - %s", broker.OrderPaid))
			defer messageSpan.End()

			//headers := broker.InjectAMQPHeaders(amqpContext)

			// publish a message
			h.channel.PublishWithContext(amqpContext, broker.OrderPaid, "", false, false, amqp.Publishing{
				ContentType:  "application/json",
				Body:         marshalledOrder,
				DeliveryMode: amqp.Persistent,
				//Headers:      headers,
			})
			log.Printf("Message published to %s, body: %s\n", broker.OrderPaid, string(marshalledOrder))
		}
	}
	w.WriteHeader(http.StatusOK)
}
