package main

import (
	"context"
	"github.com/Nicknamezz00/gorder/pkg/middleware"
	"log"
	"net/http"
	"time"

	"github.com/Nicknamezz00/gorder/gateway/entry"
	"github.com/Nicknamezz00/gorder/pkg/discovery"
	"github.com/Nicknamezz00/gorder/pkg/discovery/consul"
	"github.com/Nicknamezz00/gorder/pkg/envutil"
	_ "github.com/joho/godotenv/autoload"
)

const (
	serviceName = "gateway"
)

var (
	// expose http port to the outside
	httpAddr   = envutil.EnvString("HTTP_ADDR", ":8080")
	consulAddr = envutil.EnvString("CONSUL_ADDR", "127.0.0.1:8500")
	jaegerAddr = envutil.EnvString("JAEGER_ADDR", "127.0.0.1:4318")
)

func main() {
	shutdownTracerProvider, err := middleware.SetGlobalTracer(context.Background(), serviceName, jaegerAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdownTracerProvider(context.Background()); err != nil {
			log.Fatalf("failed to shutdown TracerProvider: %s", err)
		}
	}()

	registry, err := consul.NewRegistry(consulAddr, serviceName)
	if err != nil {
		log.Fatal(err)
	}

	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(context.Background(), instanceID, serviceName, httpAddr); err != nil {
		log.Fatal(err)
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

	mux := http.NewServeMux()
	ordersEntry := entry.NewGRPCEntry(registry)
	handler := NewHandler(ordersEntry)
	handler.registerRoutes(mux)

	log.Printf("starting http server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start http server")
	}

}
