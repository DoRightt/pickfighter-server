package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"pickfighter.com/gen"
	"pickfighter.com/internal/grpcutil"
	"pickfighter.com/pickfighter/pkg/model"
)

// ServiceHealthCheck connects to the fighters-service via gRPC to check its health status.
// It creates a new client for the fighters-service, sends a HealthCheck request, and retrieves the response.
// If successful, it converts the response from protobuf to the internal HealthStatus model.
// Returns the health status or an error if the connection or request fails.
func (g *Gateway) ServiceHealthCheck() (*model.HealthStatus, error) {
	ctx := context.Background()
	conn, err := grpcutil.ServiceConnection(ctx, "fighters-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewFightersServiceClient(conn)

	emptyRequest := &emptypb.Empty{}
	status, err := client.HealthCheck(ctx, emptyRequest)
	if err != nil {
		return nil, err
	}

	return model.HealthStatusFromProto(status), nil
}
