package entry

import (
	"context"
	pb "github.com/Nicknamezz00/pkg/api"
	"github.com/Nicknamezz00/pkg/discovery"
	"log"
)

type entry struct {
	registry discovery.Registry
}

func NewGRPCEntry(registry discovery.Registry) *entry {
	return &entry{
		registry: registry,
	}
}

func (e *entry) UpdateOrderAfterPaid(ctx context.Context, orderID, paymentLink string) error {
	conn, err := discovery.GetService(context.Background(), "orders", e.registry)
	if err != nil {
		log.Fatalf("failed to dial order server, err: %v", err)
	}
	defer conn.Close()

	c := pb.NewOrderServiceClient(conn)
	_, err = c.UpdateOrder(ctx, &pb.Order{
		ID:          orderID,
		Status:      "waiting_for_payment",
		PaymentLink: paymentLink,
	})
	return err
}
