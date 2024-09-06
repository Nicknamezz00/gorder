package adapters

import (
	"context"
	"github.com/Nicknamezz00/gorder/internal/common/errcode"
	pb "github.com/Nicknamezz00/gorder/internal/common/genproto/api"
	"log"
	"sync"
)

type MemoryOrderRepository struct {
	inMemStore []*pb.Order
	lock       *sync.RWMutex
}

func (m MemoryOrderRepository) Create(ctx context.Context, req *pb.CreateOrderRequest, items []*pb.Item) (string, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	id := "37"
	m.inMemStore = append(m.inMemStore, &pb.Order{
		ID:          id,
		CustomerID:  req.CustomerID,
		Status:      "pending",
		Items:       items,
		PaymentLink: "",
	})
	return id, nil
}

func (m MemoryOrderRepository) Get(ctx context.Context, orderID, customerID string) (*pb.Order, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, o := range m.inMemStore {
		if o.ID == orderID && o.CustomerID == customerID {
			return o, nil
		}
	}
	return nil, errcode.ErrOrderNotFound
}

func (m MemoryOrderRepository) Update(ctx context.Context, orderID string, o *pb.Order) error {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for i, v := range m.inMemStore {
		if v.ID == orderID {
			m.inMemStore[i].Status = o.Status
			m.inMemStore[i].PaymentLink = o.PaymentLink
			log.Printf("Order %s Updated! new status: %s", orderID, o.Status)
			return nil
		}
	}
	return errcode.ErrOrderNotFound
}
