package grpc

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"
	"pickfighter.com/auth/pkg/model"
	"pickfighter.com/gen"
)

// HealthCheck handles the gRPC request for checking the health status of the application.
// It delegates the call to the controller and converts the result into a protobuf response.
func (h *Handler) HealthCheck(ctx context.Context, req *emptypb.Empty) (*gen.HealthResponse, error) {
	status := h.ctrl.HealthCheck()

	return model.HealthStatusToProto(status), nil
}
