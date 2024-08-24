package grpc

import (
	"context"
	"errors"

	"pickfighter.com/fighters/internal/controller/fighters"
	"pickfighter.com/fighters/pkg/model"
	"pickfighter.com/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FightersController interface {
	SearchFightersCount(ctx context.Context, req *model.FightersRequest) (int32, error)
	SearchFighters(ctx context.Context, req *model.FightersRequest) ([]*model.Fighter, error)
}

// Handler defines a Fighters gRPC handler.
type Handler struct {
	gen.UnimplementedFightersServiceServer
	ctrl FightersController
}

// New creates a new Fighters gRPC handler.
func New(ctrl FightersController) *Handler {
	return &Handler{ctrl: ctrl}
}

// SearchFightersCount retrieves the count of fighters based on the provided request.
// It converts the request to the internal model, calls the controller's method, and returns the count.
func (h *Handler) SearchFightersCount(ctx context.Context, req *gen.FightersRequest) (*gen.FightersCountResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	fReq := &model.FightersRequest{
		Status: req.Status,
	}

	v, err := h.ctrl.SearchFightersCount(ctx, fReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.FightersCountResponse{Count: v}, nil
}

// SearchFighters retrieves fighters based on the provided request.
// It converts the request to the internal model, calls the controller's method, and returns the response.
// If no fighters are found, it returns a NotFound error; otherwise, it returns the list of fighters.
func (h *Handler) SearchFighters(ctx context.Context, req *gen.FightersRequest) (*gen.FightersResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	fReq := &model.FightersRequest{
		Status:      req.Status,
		FightersIds: req.FightersIds,
	}

	f, err := h.ctrl.SearchFighters(ctx, fReq)
	if err != nil && errors.Is(err, fighters.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.FightersResponse{
		Fighters: model.FightersToProto(f),
	}, nil
}
