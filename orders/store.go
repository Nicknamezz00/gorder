package main

import (
	"context"
	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/errcode"
	"log"
)

type OrderStore interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, orderID, customerID string) (*pb.Order, error)
	Update(ctx context.Context, orderID string, o *pb.Order) error
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
		ID:          id,
		CustomerID:  req.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
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

func (s *store) Update(ctx context.Context, orderID string, o *pb.Order) error {
	for i, v := range inMemoryStore {
		if v.ID == orderID {
			inMemoryStore[i].Status = o.Status
			inMemoryStore[i].PaymentLink = o.PaymentLink
			log.Printf("Order %s Updated! new status: %s", orderID, o.Status)
			return nil
		}
	}
	return errcode.ErrOrderNotFound
}
