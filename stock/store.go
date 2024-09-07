package main

import (
	"context"
	"fmt"

	pb "github.com/Nicknamezz00/gorder/pkg/api"
)

type StockStore interface {
	GetItem(ctx context.Context, id string) (*pb.Item, error)
	GetItems(ctx context.Context, ids []string) ([]*pb.Item, error)
}

type Store struct {
	stock map[string]*pb.Item
}

func NewStore() *Store {
	return &Store{
		stock: map[string]*pb.Item{
			"1": {
				ID:       "1",
				Name:     "Cheese Burger",
				PriceID:  "price_1PZEDmRuyMJmUCSsNZPk8lJF",
				Quantity: 20,
			},
		},
	}
}

func (s *Store) GetItem(ctx context.Context, id string) (*pb.Item, error) {
	for _, item := range s.stock {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, fmt.Errorf("item not found")
}

func (s *Store) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	var res []*pb.Item
	for _, id := range ids {
		if i, ok := s.stock[id]; ok {
			res = append(res, i)
		}
	}
	return res, nil
}
