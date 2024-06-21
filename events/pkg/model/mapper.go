package model

import (
	"fightbettr.com/gen"
)

func EventRequestFromProto(p *gen.CreateEventRequest) *EventsRequest {
	return &EventsRequest{
		Name:   p.Name,
		Fights: FightsFromProto(p.Fights),
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

func EventsToProto(events []Event) []*gen.Event {
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

func BetsToProto(bets []Bet) []*gen.Bet {
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
