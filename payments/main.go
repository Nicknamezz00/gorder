package main

import (
	"context"
	"log"
	"net"
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
	grpcAddr     = envutil.EnvString("GRPC_ADDR", "127.0.0.1:5001")
	consulAddr   = envutil.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	amqpUser     = envutil.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = envutil.EnvString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = envutil.EnvString("RABBITMQ_HOST", "127.0.0.1")
	amqpPort     = envutil.EnvString("RABBITMQ_PORT", "5672")
	stripeKey    = envutil.EnvString("STRIPE_KEY", "")
)

func main() {
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

	// stripe hook
	stripe.Key = stripeKey
	stripeProcessor := stripeProcessor.NewProcessor()

	// broker
	ch, connClose := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		ch.Close()
		connClose()
	}()

	// service
	svc := NewService(stripeProcessor)
	// rabbitmq
	amqpConsumer := NewConsumer(svc)
	go amqpConsumer.Listen(ch)

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
