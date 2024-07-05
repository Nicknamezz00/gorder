package main

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/broker"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service   OrderService
	mqChannel *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service OrderService, mqChannel *amqp.Channel) {
	h := &grpcHandler{
		service:   service,
		mqChannel: mqChannel,
	}
	pb.RegisterOrderServiceServer(grpcServer, h)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("new order created, order: %v", req)
	o := &pb.Order{
		ID: "37",
	}
	q, err := h.mqChannel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		log.Fatal(err)
	}
	h.mqChannel.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Body:         marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})
	return o, nil
}
