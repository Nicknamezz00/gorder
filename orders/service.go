package main

import (
	"context"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"github.com/Nicknamezz00/gorder/pkg/errcode"
	"log"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest, []*pb.Item) (*pb.Order, error)
	ValidateOrder(context.Context, *pb.CreateOrderRequest) ([]*pb.Item, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
	UpdateOrder(context.Context, *pb.Order) (*pb.Order, error)
}

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{
		store: store,
	}
}

func (s *service) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	return s.store.Get(ctx, req.OrderID, req.CustomerID)
}

func (s *service) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest, items []*pb.Item) (*pb.Order, error) {
	id, err := s.store.Create(ctx, req, items)
	if err != nil {
		return nil, err
	}
	o := &pb.Order{
		ID:         id,
		CustomerID: req.CustomerID,
		Status:     "pending",
		Items:      items,
	}
	return o, nil
}

func (s *service) UpdateOrder(ctx context.Context, o *pb.Order) (*pb.Order, error) {
	err := s.store.Update(ctx, o.ID, o)
	if err != nil {
		return nil, err
	}
	return o, nil
}

func (s *service) ValidateOrder(ctx context.Context, req *pb.CreateOrderRequest) ([]*pb.Item, error) {
	if len(req.Items) == 0 {
		return nil, errcode.ErrNoItems
	}
	var itemsWithPrice []*pb.Item
	mergedItems := packItems(req.Items)
	log.Printf("merged items: %v", mergedItems)
	for _, it := range mergedItems {
		// PRICEID
		itemsWithPrice = append(itemsWithPrice, &pb.Item{
			PriceID:  "price_1PZEDmRuyMJmUCSsNZPk8lJF",
			ID:       it.ID,
			Quantity: it.Quantity,
		})
	}
	return itemsWithPrice, nil
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
