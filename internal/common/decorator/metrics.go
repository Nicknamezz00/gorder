package decorator

import "context"

type MetricsClient struct {
	// Implement this;
}

type queryMetricsDecorator[C any, R any] struct {
	base   QueryHandler[C, R]
	client MetricsClient
}

func (q queryMetricsDecorator[C, R]) Handle(ctx context.Context, cmd C) (R, error) {
	panic("implement me")
}

type commandMetricsDecorator[C any] struct {
	base   CommandHandler[C]
	client MetricsClient
}

func (c commandMetricsDecorator[C]) Handle(ctx context.Context, cmd C) error {
	panic("implement me")
}

