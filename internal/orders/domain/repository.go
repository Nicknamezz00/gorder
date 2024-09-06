package order

import (
	"context"
	pb "github.com/Nicknamezz00/gorder/internal/common/genproto/api"
)

type Repository interface {
	Create(context.Context, *pb.CreateOrderRequest, []*pb.Item) (string, error)
	Get(ctx context.Context, orderID, customerID string) (*pb.Order, error)
	Update(ctx context.Context, orderID string, o *pb.Order) error
}

