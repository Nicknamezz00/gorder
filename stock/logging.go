package main

import (
	"context"
	pb "github.com/Nicknamezz00/gorder/pkg/api"
	"go.uber.org/zap"
	"time"
)

type LoggingMiddleware struct {
	next StockService
}

func (s *LoggingMiddleware) CheckIfItemAreInStock(ctx context.Context, req []*pb.ItemWithQuantity) (bool, []*pb.Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("CeckIfItemAreInStock", zap.Duration("took", time.Since(start)))
	}()
	return s.next.CheckIfItemAreInStock(ctx, req)
}

func (s *LoggingMiddleware) GetItems(ctx context.Context, ids []string) ([]*pb.Item, error) {
	start := time.Now()
	defer func() {
		zap.L().Info("GetItems", zap.Duration("took", time.Since(start)))
	}()
	return s.next.GetItems(ctx, ids)
}

func NewLoggingMiddleware(next StockService) StockService {
	return &LoggingMiddleware{next}
}
