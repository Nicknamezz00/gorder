package main

import (
	"context"
	"testing"

	pb "github.com/Nicknamezz00/pkg/api"
)

type testProcessor struct{}

func NewTestProcessor() *testProcessor {
	return &testProcessor{}
}

func (t *testProcessor) CreatePaymentLink(*pb.Order) (string, error) {
	return "test-link", nil
}

func TestService(t *testing.T) {
	p := NewTestProcessor()
	svc := NewService(p)
	t.Run("create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &pb.Order{})
		if err != nil {
			t.Errorf("create payment error, want nil, got: %v", err)
		}
		if link == "" {
			t.Error("want payment link, got empty link")
		}
	})
}
