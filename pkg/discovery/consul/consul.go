package consul

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	capi "github.com/hashicorp/consul/api"
)

type Registry struct {
	client *capi.Client
}

func NewRegistry(addr, serviceName string) (*Registry, error) {
	config := capi.DefaultConfig()
	config.Address = addr
	client, err := capi.NewClient(config)
	if err != nil {
		return nil, err
	}

	return &Registry{
		client: client,
	}, nil
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("invalid host:port format")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	host := parts[0]
	return r.client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      instanceID,
		Address: host,
		Port:    port,
		Name:    serviceName,
		Check: &capi.AgentServiceCheck{
			CheckID:                        instanceID,
			TLSSkipVerify:                  true,
			TTL:                            "5s",
			Timeout:                        "1s",
			DeregisterCriticalServiceAfter: "10s",
		},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceID, serviceName string) error {
	log.Printf("deregistering service %s, instanceID: %s", serviceName, instanceID)
	return r.client.Agent().CheckDeregister(instanceID)
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := r.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		ids = append(ids, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return ids, nil
}

func (r *Registry) HeartBeat(instanceID, serviceName string) error {
	return r.client.Agent().UpdateTTL(instanceID, "online", api.HealthPassing)
}
