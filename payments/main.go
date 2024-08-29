package main

import (
	"context"
	"github.com/Nicknamezz00/gorder-payments/entry"
	"github.com/Nicknamezz00/pkg/middleware"
	"log"
	"net"
	"net/http"
	"time"

	stripeProcessor "github.com/Nicknamezz00/gorder-payments/processor/stripe"
	"github.com/Nicknamezz00/pkg/broker"
	"github.com/Nicknamezz00/pkg/discovery"
	"github.com/Nicknamezz00/pkg/discovery/consul"
	"github.com/Nicknamezz00/pkg/envutil"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stripe/stripe-go/v79"
	"google.golang.org/grpc"
)

const (
	serviceName = "payments"
)

var (
	// expose grpc port to the outside
	httpAddr     = envutil.EnvString("HTTP_ADDR", "127.0.0.1:8081")
	grpcAddr     = envutil.EnvString("GRPC_ADDR", "127.0.0.1:5001")
	consulAddr   = envutil.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	amqpUser     = envutil.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = envutil.EnvString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = envutil.EnvString("RABBITMQ_HOST", "127.0.0.1")
	amqpPort     = envutil.EnvString("RABBITMQ_PORT", "5672")
	stripeKey    = envutil.EnvString("STRIPE_KEY", "")
	jaegerAddr   = envutil.EnvString("JAEGER_ADDR", "127.0.0.1:4318")
	//endpointStripeSecret: stripe listen --forward-to localhost:8081/webhook
	endpointStripeSecret = envutil.EnvString("ENDPOINT_STRIPE_SECRET", "whsec....")
)

func main() {
	if err := middleware.SetGlobalTracer(context.TODO(), serviceName, jaegerAddr); err != nil {
		log.Fatal("failed to set global tracer")
	}

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		panic(err)
	}
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(context.Background(), instanceID, serviceName, grpcAddr); err != nil {
		panic(err)
	}
	go func() {
		for {
			if err := registry.HeartBeat(instanceID, serviceName); err != nil {
				log.Fatalf("no heartbeat: %s", serviceName)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(context.Background(), instanceID, serviceName)

	// stripeProcessor hook
	stripe.Key = stripeKey
	stripeProc := stripeProcessor.NewProcessor()

	// broker
	ch, connClose := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		ch.Close()
		connClose()
	}()

	// service
	paymentGateway := entry.NewGRPCEntry(registry)
	svc := NewService(stripeProc, paymentGateway)

	// rabbitmq
	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

	// http server
	mux := http.NewServeMux()
	httpServer := NewPaymentHandler(ch)
	httpServer.registerRoutes(mux)
	go func() {
		log.Printf("starting payment http server at %s", httpAddr)
		if err := http.ListenAndServe(httpAddr, mux); err != nil {
			log.Fatalf("failed to start payment http server, err: %v", err)
		}
	}()
	// grpc
	grpcSrv := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen grpc server: %v", err)
	}
	defer l.Close()

	log.Printf("starting grpc server at %s", grpcAddr)
	if err := grpcSrv.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
