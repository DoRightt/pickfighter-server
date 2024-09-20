package model

import (
	eventmodel "pickfighter.com/events/pkg/model"
	fightersmodel "pickfighter.com/fighters/pkg/model"
	"pickfighter.com/gen"
)

func HealthStatusFromProto(req *gen.HealthResponse) *HealthStatus {
	return &HealthStatus{
		AppDevVersion: req.AppDevVersion,
		AppName:       req.AppName,
		AppRunDate:    req.AppRunDate,
		AppTimeAlive:  req.AppTimeAlive,
		Healthy:       req.Healthy,
		Message:       req.Message,
		Timestamp:     req.Timestamp,
	}
}

func ServiceEventToGatewayEvent(event *eventmodel.Event, fightersList map[int32]*fightersmodel.Fighter) *Event {
	fights := event.Fights
	updatedEvent := &Event{
		EventId: event.EventId,
		Name:    event.Name,
		IsDone:  event.IsDone,
		Fights:  make([]Fight, len(fights)),
	}

	for i, v := range fights {
		redId := v.FighterRedId
		blueId := v.FighterBlueId

		fight := Fight{
			FightId:     v.FightId,
			EventId:     v.EventId,
			FighterRed:  *fightersList[redId],
			FighterBlue: *fightersList[blueId],
			IsDone:      v.IsDone,
			IsCanceled:  v.IsCanceled,
			NotContest:  v.NotContest,
			Result:      v.Result,
			CreatedAt:   v.CreatedAt,
			FightDate:   v.FightDate,
		}

		updatedEvent.Fights[i] = fight

	}
	return updatedEvent
}
