package main

import "context"

type OrderStore interface {
	Create(context.Context) error
}
type store struct {
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(context.Context) error {
	return nil
}
