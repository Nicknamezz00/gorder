package main

import (
	"context"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

type StockGRPCHandler struct {
	pb.UnimplementedStockServiceServer

	service StockService
	channel *amqp.Channel
}

func NewGRPCHandler(
	server *grpc.Server,
	stockService StockService,
	channel *amqp.Channel,
) {
	handler := &StockGRPCHandler{
		service: stockService,
		channel: channel,
	}

	pb.RegisterStockServiceServer(server, handler)
}

func (s *StockGRPCHandler) CheckIfItemIsInStock(ctx context.Context, p *pb.CheckIfItemIsInStockRequest) (*pb.CheckIfItemIsInStockResponse, error) {
	inStock, items, err := s.service.CheckIfItemAreInStock(ctx, p.Items)
	if err != nil {
		return nil, err
	}
	return &pb.CheckIfItemIsInStockResponse{
		InStock: inStock,
		Items:   items,
	}, nil
}

func (s *StockGRPCHandler) GetItems(ctx context.Context, payload *pb.GetItemsRequest) (*pb.GetItemsResponse, error) {
	items, err := s.service.GetItems(ctx, payload.ItemIDs)
	if err != nil {
		return nil, err
	}
	return &pb.GetItemsResponse{
		Items: items,
	}, nil
}
