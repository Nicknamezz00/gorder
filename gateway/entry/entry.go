package entry

import (
	"context"

	pb "github.com/Nicknamezz00/pkg/api"
)

type OrdersEntry interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, orderID, customerID string) (*pb.Order, error)
}
