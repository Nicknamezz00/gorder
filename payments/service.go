package main

import (
	"context"
	"github.com/Nicknamezz00/gorder/payments/entry"

	"github.com/Nicknamezz00/gorder/payments/processor"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}

type Service struct {
	processor processor.PaymentProcessor
	entry     entry.PaymentEntry
}

func NewService(processor processor.PaymentProcessor, paymentEntry entry.PaymentEntry) *Service {
	return &Service{
		processor: processor,
		entry:     paymentEntry,
	}
}

func (s *Service) CreatePayment(ctx context.Context, o *pb.Order) (string, error) {
	link, err := s.processor.CreatePaymentLink(o)
	if err != nil {
		return "", err
	}
	err = s.entry.UpdateOrderAfterPaid(ctx, o.ID, link)
	if err != nil {
		return "", err
	}
	return link, nil
}
