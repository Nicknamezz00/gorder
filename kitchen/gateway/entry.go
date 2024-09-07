package gateway

import (
	"context"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
)

type KitchenGateway interface {
	UpdateOrder(context.Context, *pb.Order) error
}
