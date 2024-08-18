package processor

import (
	pb "github.com/Nicknamezz00/pkg/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(*pb.Order) (string, error)
}
