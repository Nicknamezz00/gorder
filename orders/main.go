package main

import (
	"context"
	"log"
	"net"

	common "github.com/Nicknamezz00/gorder-common"
	"google.golang.org/grpc"
)

const (
	serviceName = "orders"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "127.0.0.1:5000")
)

func main() {
	grpcSrv := grpc.NewServer()
	l, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		log.Fatalf("failed to listen grpc server: %v", err)
	}
	defer l.Close()

	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcSrv, svc)

	log.Printf("starting grpc server at %s", grpcAddr)
	svc.CreateOrder(context.Background())
	if err := grpcSrv.Serve(l); err != nil {
		log.Fatal(err.Error())
	}
}
