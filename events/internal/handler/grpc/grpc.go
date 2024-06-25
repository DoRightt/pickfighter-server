package grpc

import (
	"context"

	"fightbettr.com/events/internal/controller/event"
	"fightbettr.com/events/pkg/model"
	"fightbettr.com/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler defines a Event gRPC handler.
type Handler struct {
	gen.UnimplementedEventServiceServer
	ctrl *event.Controller
}

// New creates a new Event gRPC handler.
func New(ctrl *event.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) CreateEvent(ctx context.Context, req *gen.CreateEventRequest) (*gen.CreateEventResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	eventReq := model.EventRequestFromProto(req)
	v, err := h.ctrl.CreateEvent(ctx, eventReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.CreateEventResponse{EventId: v}, nil
}

func (h *Handler) GetEvents(ctx context.Context, req *gen.GetEventsRequest) (*gen.GetEventsResponse, error) {
	resp, err := h.ctrl.GetEvents(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	events := model.EventsToProto(resp.Events)

	return &gen.GetEventsResponse{Count: resp.Count, Events: events}, nil
}

func (h *Handler) CreateBet(ctx context.Context, req *gen.CreateBetRequest) (*gen.CreateBetResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	betReq := model.BetRequestFromProto(req)
	v, err := h.ctrl.CreateBet(ctx, betReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.CreateBetResponse{BetId: v}, nil
}

func (h *Handler) GetBets(ctx context.Context, req *gen.BetsRequest) (*gen.BetsResponse, error) {
	resp, err := h.ctrl.GetBets(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	bets := model.BetsToProto(resp.Bets)

	return &gen.BetsResponse{Bets: bets, Count: resp.Count}, nil
}

func (h *Handler) SetResult(ctx context.Context, req *gen.FightResultRequest) (*gen.FightResultResponse, error) {
	if req == nil {
		return nil, status.Errorf(codes.InvalidArgument, "nil request")
	}

	resReq := model.FightResultFromProto(req)
	_, err := h.ctrl.SetFightResult(ctx, resReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &gen.FightResultResponse{}, nil
}

// GracefulShutdown initiates a graceful shutdown of the service by delegating the signal handling to the controller.
func (h *Handler) GracefulShutdown(ctx context.Context, sig string) {
	h.ctrl.GracefulShutdown(ctx, sig)
}
