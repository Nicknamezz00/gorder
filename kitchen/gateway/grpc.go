package gateway

import (
	"context"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"github.com/Nicknamezz00/gorder/pkg/discovery"
	"google.golang.org/grpc"
	"log"
)

type Gateway struct {
	registry discovery.Registry
}

func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) UpdateOrder(ctx context.Context, o *pb.Order) error {
	conn, err := discovery.GetService(context.Background(), "orders", g.registry)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	ordersClient := pb.NewOrderServiceClient(conn)
	_, err = ordersClient.UpdateOrder(ctx, o)
	return err
}
