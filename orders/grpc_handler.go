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

func (h *grpcHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, req)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	items, err := h.service.ValidateOrder(ctx, req)
	if err != nil {
		return nil, err
	}
	o, err := h.service.CreateOrder(ctx, req, items)
	if err != nil {
		return nil, err
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
