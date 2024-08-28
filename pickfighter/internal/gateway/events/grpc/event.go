package grpc

import (
	"context"

	eventmodel "pickfighter.com/events/pkg/model"
	"pickfighter.com/gen"
	"pickfighter.com/internal/grpcutil"
	"pickfighter.com/pkg/discovery"
)

// Gateway defines an gRPC gateway for a event service.
type Gateway struct {
	registry discovery.Registry
}

// New creates a new gRPC gateway for a event service.
func New(registry discovery.Registry) *Gateway {
	return &Gateway{registry}
}

func (g *Gateway) CreateEvent(ctx context.Context, req *eventmodel.EventRequest) (*eventmodel.Event, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "event-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewEventServiceClient(conn)

	eventReq := eventmodel.EventRequestToProto(req)
	resp, err := client.CreateEvent(ctx, eventReq)
	if err != nil {
		return nil, err
	}

	event := &eventmodel.Event{EventId: resp.EventId}

	return event, nil
}

func (g *Gateway) SearchEvents(ctx context.Context) (*eventmodel.EventsResponse, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "event-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewEventServiceClient(conn)

	resp, err := client.GetEvents(ctx, &gen.GetEventsRequest{}) // TODO empty request
	if err != nil {
		return nil, err
	}

	events := &eventmodel.EventsResponse{Count: resp.Count, Events: eventmodel.EventsFromProto(resp.Events)}

	return events, nil
}

func (g *Gateway) CreateBet(ctx context.Context, req *eventmodel.Bet) (*eventmodel.Bet, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "event-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewEventServiceClient(conn)

	betReq := eventmodel.BetRequestToProto(req)
	resp, err := client.CreateBet(ctx, betReq)
	if err != nil {
		return nil, err
	}

	bet := &eventmodel.Bet{BetId: resp.BetId}

	return bet, nil
}

func (g *Gateway) SearchBets(ctx context.Context, userId int32) (*eventmodel.BetsResponse, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "event-service", g.registry)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := gen.NewEventServiceClient(conn)

	resp, err := client.GetBets(ctx, &gen.BetsRequest{UserId: userId}) // TODO empty request
	if err != nil {
		return nil, err
	}

	bets := &eventmodel.BetsResponse{Count: resp.Count, Bets: eventmodel.BetsFromProto(resp.Bets)}

	return bets, nil
}

func (g *Gateway) SetResult(ctx context.Context, req *eventmodel.FightResultRequest) (int32, error) {
	conn, err := grpcutil.ServiceConnection(ctx, "event-service", g.registry)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := gen.NewEventServiceClient(conn)

	resultReq := eventmodel.FightResultToProto(req)
	resp, err := client.SetResult(ctx, resultReq)
	if err != nil {
		return 0, err
	}

	return resp.FightId, nil
}
