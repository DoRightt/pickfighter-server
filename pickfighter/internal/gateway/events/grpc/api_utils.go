package grpc

import (
	"context"

	"github.com/DoRightt/pickfighter-server/gen"
	"github.com/DoRightt/pickfighter-server/internal/grpcutil"
	"github.com/DoRightt/pickfighter-server/pickfighter/pkg/model"
	"google.golang.org/protobuf/types/known/emptypb"
)

// ServiceHealthCheck connects to the events-service via gRPC to check its health status.
// It creates a new client for the events-service, sends a HealthCheck request, and retrieves the response.
// If successful, it converts the response from protobuf to the internal HealthStatus model.
// Returns the health status or an error if the connection or request fails.
func (g *Gateway) ServiceHealthCheck() (*model.HealthStatus, error) {
	ctx := context.Background()
	conn, err := grpcutil.ServiceConnection(ctx, "events-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewEventServiceClient(conn)

	emptyRequest := &emptypb.Empty{}
	status, err := client.HealthCheck(ctx, emptyRequest)
	if err != nil {
		return nil, err
	}

	return model.HealthStatusFromProto(status), nil
}
