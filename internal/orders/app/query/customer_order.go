package query

import (
	"context"
	"github.com/Nicknamezz00/gorder/internal/common/decorator"
	pb "github.com/Nicknamezz00/gorder/internal/common/genproto/api"
	order "github.com/Nicknamezz00/gorder/internal/orders/domain"
	"go.uber.org/zap"
)

type CustomerOrder struct {
	orderID    string
	customerID string
}

type customerOrderHandler struct {
	orderRepo order.Repository
}

type CustomerOrderHandler decorator.QueryHandler[CustomerOrder, *pb.Order]

func NewCustomerOrderHandler(
	orderRepo order.Repository,
	logger *zap.Logger,
	metricsClient decorator.MetricsClient,
) CustomerOrderHandler {
	return decorator.ApplyQueryDecorators[CustomerOrder, *pb.Order](
		customerOrderHandler{orderRepo: orderRepo},
		logger,
		metricsClient,
	)
}

func (c customerOrderHandler) Handle(ctx context.Context, query CustomerOrder) (*pb.Order, error) {
	o, err := c.orderRepo.Get(ctx, query.orderID, query.customerID)
	if err != nil {
		return nil, err
	}
	return o, nil
}

