package model

import (
	"pickfighter.com/gen"
)

func EventRequestFromProto(p *gen.CreateEventRequest) *EventRequest {
	return &EventRequest{
		Name:   p.Name,
		Fights: FightsFromProto(p.Fights),
	}
}

func EventRequestToProto(req *EventRequest) *gen.CreateEventRequest {
	return &gen.CreateEventRequest{
		Name:   req.Name,
		Fights: FightsToProto(req.Fights),
	}
}

func FightsFromProto(p []*gen.Fight) []Fight {
	fights := make([]Fight, len(p))

	for i, v := range p {
		fights[i] = Fight{
			FightId:       v.FightId,
			EventId:       v.EventId,
			FighterRedId:  v.FighterRedId,
			FighterBlueId: v.FighterBlueId,
			IsDone:        v.IsDone,
			IsCanceled:    v.IsCanceled,
			Result:        v.Result,
			CreatedAt:     v.CreatedAt,
			FightDate:     int(v.FightDate),
		}
	}

	return fights
}

func FightsToProto(fights []Fight) []*gen.Fight {
	protoFights := make([]*gen.Fight, len(fights))

	for i, v := range fights {
		protoFights[i] = &gen.Fight{
			FightId:       v.FightId,
			EventId:       v.EventId,
			FighterRedId:  v.FighterRedId,
			FighterBlueId: v.FighterBlueId,
			IsDone:        v.IsDone,
			IsCanceled:    v.IsCanceled,
			Result:        v.Result,
			CreatedAt:     v.CreatedAt,
			FightDate:     int64(v.FightDate),
		}
	}

	return protoFights
}

func EventsFromProto(p []*gen.Event) []*Event {
	events := make([]*Event, len(p))

	for i, v := range p {
		events[i] = &Event{
			EventId: v.EventId,
			Name:    v.Name,
			IsDone:  v.IsDone,
			Fights:  FightsFromProto(v.Fights),
		}
	}

	return events
}

func EventsToProto(events []*Event) []*gen.Event {
	protoEvents := make([]*gen.Event, len(events))

	for i, v := range events {
		protoEvents[i] = &gen.Event{
			EventId: v.EventId,
			Name:    v.Name,
			IsDone:  v.IsDone,
			Fights:  FightsToProto(v.Fights),
		}
	}

	return protoEvents
}

func BetRequestFromProto(p *gen.CreateBetRequest) *Bet {
	return &Bet{
		BetId:     p.BetId,
		FightId:   p.FightId,
		UserId:    p.UserId,
		FighterId: p.FighterId,
	}
}

func BetRequestToProto(bet *Bet) *gen.CreateBetRequest {
	return &gen.CreateBetRequest{
		BetId:     bet.BetId,
		FightId:   bet.FightId,
		UserId:    bet.UserId,
		FighterId: bet.FighterId,
	}
}

func BetsFromProto(p []*gen.Bet) []*Bet {
	bets := make([]*Bet, len(p))

	for i, v := range p {
		bets[i] = &Bet{
			BetId:     v.BetId,
			FightId:   v.FightId,
			UserId:    v.UserId,
			FighterId: v.FighterId,
		}
	}

	return bets
}

func BetsToProto(bets []*Bet) []*gen.Bet {
	protoBets := make([]*gen.Bet, len(bets))

	for i, v := range bets {
		protoBets[i] = &gen.Bet{
			BetId:     v.BetId,
			FightId:   v.FightId,
			UserId:    v.UserId,
			FighterId: v.FighterId,
		}
	}

	return protoBets
}

func FightResultFromProto(p *gen.FightResultRequest) *FightResultRequest {
	return &FightResultRequest{
		FightId:    p.FightId,
		WinnerId:   p.WinnerId,
		NotContest: p.NotContest,
	}
}

func FightResultToProto(req *FightResultRequest) *gen.FightResultRequest {
	return &gen.FightResultRequest{
		FightId:    req.FightId,
		WinnerId:   req.WinnerId,
		NotContest: req.NotContest,
	}
}
