package decorator

import (
	"context"
	"go.uber.org/zap"
)

// QueryHandler receives a query of type Q, returns a result of type R.
type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, query Q) (R, error)
}

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], logger *zap.Logger, metricsClient MetricsClient) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: queryMetricsDecorator[H, R]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}

