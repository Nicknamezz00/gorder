package main

import (
	"context"
	"github.com/Nicknamezz00/pkg/middleware"
	"log"
	"net"
	"time"

	"github.com/Nicknamezz00/pkg/broker"
	"github.com/Nicknamezz00/pkg/discovery"
	"github.com/Nicknamezz00/pkg/discovery/consul"
	"github.com/Nicknamezz00/pkg/envutil"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

const (
	serviceName = "orders"
)

var (
	// expose grpc port to the outside
	grpcAddr     = envutil.EnvString("GRPC_ADDR", "127.0.0.1:5000")
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

	if err := middleware.SetGlobalTracer(context.Background(), serviceName, jaegerAddr); err != nil {
		logger.Fatal("could set global tracer", zap.Error(err))
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
	NewGRPCHandler(grpcSrv, svc, ch)

	logger.Info("starting grpc server at %s", zap.String("grpcAddr", grpcAddr))
	// svc.CreateOrder(context.Background())
	if err := grpcSrv.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
