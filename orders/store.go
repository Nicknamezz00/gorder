package main

import (
	"context"
	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/errcode"
)

type OrderStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, orderID, customerID string) (*pb.Order, error)
}

type store struct {
}

var inMemoryStore = make([]*pb.Order, 0)

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, req *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	id := "37"
	inMemoryStore = append(inMemoryStore, &pb.Order{
		ID:         id,
		CustomerID: req.CustomerID,
		Status:     "pending",
		Items:      items,
	})
	return id, nil
}

func (s *store) Get(ctx context.Context, orderID, customerID string) (*pb.Order, error) {
	for _, o := range inMemoryStore {
		if o.ID == orderID && o.CustomerID == customerID {
			return o, nil
		}
	}
	return nil, errcode.ErrOrderNotFound
}
