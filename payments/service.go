package main

import (
	"context"

	"github.com/Nicknamezz00/gorder-payments/processor"
	pb "github.com/Nicknamezz00/pkg/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{
		processor: processor,
	}
}

func (s *service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}
	return link, nil
}
