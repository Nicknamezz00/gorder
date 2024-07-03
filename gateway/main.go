package main

import (
	"log"
	"net/http"

	common "github.com/Nicknamezz00/gorder-common"
	pb "github.com/Nicknamezz00/gorder-common/api"
	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	serviceName = "gateway"
)

var (
	httpAddr     = common.EnvString("HTTP_ADDR", ":8080")
	orderService = "127.0.0.1:5000"
)

func main() {
	conn, err := grpc.Dial(orderService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()
	log.Printf("dialing order service at: %s", orderService)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.registerRoutes(mux)

	log.Printf("starting http server at %s", httpAddr)

	if err := http.ListenAndServe(httpAddr, mux); err != nil {
		log.Fatal("failed to start http server")
	}

}
