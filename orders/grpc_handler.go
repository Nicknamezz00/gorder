package main

import (
	"context"
	"log"

	pb "github.com/Nicknamezz00/pkg/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service OrderService
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService) {
	h := &grpcHandler{
		service: service,
	}
	pb.RegisterOrderServiceServer(grpcServer, h)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("new order created, order: %v", req)
	o := &pb.Order{
		ID: "37",
	}
	return o, nil
}
