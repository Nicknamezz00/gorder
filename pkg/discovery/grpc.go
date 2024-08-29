package discovery

import (
	"context"
	"math/rand"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetService(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
  addrs, err := registry.Discover(ctx, serviceName)
  if err != nil {
    return nil, err
  }
  i := rand.Intn(len(addrs))
  return grpc.NewClient(
    addrs[i],
    grpc.WithTransportCredentials(insecure.NewCredentials()),
    grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
    grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
  )
}
