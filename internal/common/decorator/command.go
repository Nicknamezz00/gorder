package decorator

import (
	"context"
	"go.uber.org/zap"
)

// CommandHandler receives a command of type C, returns an error.
type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func ApplyCommandDecorators[C any](handler CommandHandler[C], logger *zap.Logger, metricsClient MetricsClient) CommandHandler[C] {
	return commandLoggingDecorator[C]{
		base: commandMetricsDecorator[C]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}

