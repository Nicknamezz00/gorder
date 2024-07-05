package main

import (
	"context"
	"log"

	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/errcode"
)

type OrderService interface {
	CreateOrder(context.Context) error
	ValidateOrder(context.Context, *pb.CreateOrderRequest) error
}

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateOrder(context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) error {
	if len(req.Items) == 0 {
		return errcode.ErrNoItems
	}
	// items := packItems(req.Items)
	log.Printf("packed items: %v", packItems(req.Items))
	log.Printf("slow packed items: %v", packItemsSlow(req.Items))
	return nil
}

// packItems merges quantities of the same item.
func packItems(items []*pb.ItemWithQuantity) []*pb.ItemWithQuantity {
	merged := make(map[string]int32)
	for _, item := range items {
		merged[item.ID] += item.Quantity
	}
	result := make([]*pb.ItemWithQuantity, 0, len(merged))
	for id, quantity := range merged {
		result = append(result, &pb.ItemWithQuantity{ID: id, Quantity: quantity})
	}
	return result
}

func packItemsSlow(items []*pb.ItemWithQuantity) []*pb.ItemWithQuantity {
	merged := make([]*pb.ItemWithQuantity, 0)
	for _, item := range items {
		found := false
		for _, toUpdate := range merged {
			if item.ID == toUpdate.ID {
				toUpdate.Quantity += item.Quantity
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, item)
		}
	}
	return merged
}
