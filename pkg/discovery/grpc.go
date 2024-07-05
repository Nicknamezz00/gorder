package discovery

import (
	"context"
	"math/rand"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func GetService(ctx context.Context, serviceName string, registry Registry) (*grpc.ClientConn, error) {
	addrs, err := registry.Discover(ctx, serviceName)
	if err != nil {
		return nil, err
	}
	i := rand.Intn(len(addrs))
	return grpc.Dial(addrs[i], grpc.WithTransportCredentials(insecure.NewCredentials()))
}
