package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	ErrNotFound = errors.New("no service addresses found")
)

type Registry struct {
	client *redis.Client
}

// NewRegistry creates a new Redis-based service registry instance.
func NewRegistry(addr string) (*Registry, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &Registry{client: client}, nil
}

// Register creates a service record in the registry.
func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	err := r.client.SAdd(ctx, fmt.Sprintf("service:%s", serviceName), instanceID).Err()
	if err != nil {
		return fmt.Errorf("failed to add service instance to Redis: %w", err)
	}

	err = r.client.HSet(ctx, fmt.Sprintf("service:%s:%s", serviceName, instanceID), "address", hostPort, "status", "healthy").Err()
	if err != nil {
		return fmt.Errorf("failed to register service instance data in Redis: %w", err)
	}

	return nil
}

// Deregister removes a service record from the registry.
func (r *Registry) Deregister(ctx context.Context, instanceID string, serviceName string) error {
	err := r.client.SRem(ctx, fmt.Sprintf("service:%s", serviceName), instanceID).Err()
	if err != nil {
		return fmt.Errorf("failed to remove service instance from Redis: %w", err)
	}

	err = r.client.Del(ctx, fmt.Sprintf("service:%s:%s", serviceName, instanceID)).Err()
	if err != nil {
		return fmt.Errorf("failed to remove service instance data from Redis: %w", err)
	}

	return nil
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (r *Registry) ServiceAddresses(ctx context.Context, serviceID string) ([]string, error) {
	instanceIDs, err := r.client.SMembers(ctx, fmt.Sprintf("service:%s", serviceID)).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get service instances from Redis: %w", err)
	}

	if len(instanceIDs) == 0 {
		return nil, ErrNotFound
	}

	var addresses []string
	for _, instanceID := range instanceIDs {
		address, err := r.client.HGet(ctx, fmt.Sprintf("service:%s:%s", serviceID, instanceID), "address").Result()
		if err != nil {
			continue
		}
		addresses = append(addresses, address)
	}

	if len(addresses) == 0 {
		return nil, ErrNotFound
	}

	return addresses, nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (r *Registry) ReportHealthyState(instanceID string, serviceName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.client.HSet(ctx, fmt.Sprintf("service:%s:%s", serviceName, instanceID), "status", "healthy").Err()
	if err != nil {
		return fmt.Errorf("failed to update health status in Redis: %w", err)
	}
	return nil
}
