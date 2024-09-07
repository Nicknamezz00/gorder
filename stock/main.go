package main

import (
	"context"
	"github.com/Nicknamezz00/gorder/pkg/broker"
	"github.com/Nicknamezz00/gorder/pkg/discovery"
	"github.com/Nicknamezz00/gorder/pkg/discovery/consul"
	"github.com/Nicknamezz00/gorder/pkg/envutil"
	"github.com/Nicknamezz00/gorder/pkg/middleware"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	serviceName = "stock"
)

var (
	// expose grpc port to the outside
	grpcAddr     = envutil.EnvString("GRPC_ADDR", "127.0.0.1:5002")
	consulAddr   = envutil.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	amqpUser     = envutil.EnvString("RABBITMQ_USER", "guest")
	amqpPassword = envutil.EnvString("RABBITMQ_PASSWORD", "guest")
	amqpHost     = envutil.EnvString("RABBITMQ_HOST", "127.0.0.1")
	amqpPort     = envutil.EnvString("RABBITMQ_PORT", "5672")
	jaegerAddr   = envutil.EnvString("JAEGER_ADDR", "127.0.0.1:4318")
)

func main() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		_ = logger.Sync()
	}(logger)
	zap.ReplaceGlobals(logger)

	shutdownTracerProvider, err := middleware.SetGlobalTracer(context.Background(), serviceName, jaegerAddr)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func() {
		if err := shutdownTracerProvider(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		logger.Fatal(err.Error())
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(context.Background(), instanceID, serviceName, grpcAddr); err != nil {
		logger.Fatal(err.Error())
	}
	go func() {
		for {
			if err := registry.HeartBeat(instanceID, serviceName); err != nil {
				log.Fatalf("no heartbeat from %s to registry, err = %v", serviceName, err)
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(context.Background(), instanceID, serviceName)

	ch, connClose := broker.Connect(amqpUser, amqpPassword, amqpHost, amqpPort)
	defer func() {
		_ = ch.Close()
		_ = connClose()
	}()

	grpcSrv := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen grpc server: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	svcWithTelemetry := NewTelemetryMiddleware(svc)
	svcWithLogging := NewLoggingMiddleware(svcWithTelemetry)
	NewGRPCHandler(grpcSrv, svcWithLogging, ch)

	go NewConsumer().Listen(ch)

	logger.Info("starting grpc server at %s", zap.String("grpcAddr", grpcAddr))
	if err := grpcSrv.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
