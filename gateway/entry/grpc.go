package entry

import (
	"context"
	"log"

	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/discovery"
)

type entry struct {
	registry discovery.Registry
}

func NewGRPCEntry(registry discovery.Registry) *entry {
	return &entry{
		registry: registry,
	}
}

func (e *entry) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	conn, err := discovery.GetService(ctx, "orders", e.registry)
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
		return nil, err
	}
	c := pb.NewOrderServiceClient(conn)
	return c.CreateOrder(ctx, &pb.CreateOrderRequest{
		CustomerID: req.CustomerID,
		Items:      req.Items,
	})
}
