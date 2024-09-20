package pickfighter

import (
	"fmt"
	"time"

	eventmodel "pickfighter.com/events/pkg/model"
	fightersmodel "pickfighter.com/fighters/pkg/model"
	"pickfighter.com/pickfighter/pkg/model"
	gatewaymodel "pickfighter.com/pickfighter/pkg/model"
	"pickfighter.com/pickfighter/pkg/version"
)

func (c *Controller) getFightersIds(events []*eventmodel.Event) []int32 {
	var ids []int32
	uniqueIdsMap := make(map[int32]struct{})

	for _, event := range events {
		fights := event.Fights

		for _, fight := range fights {
			redId := fight.FighterRedId
			blueId := fight.FighterBlueId

			if _, exists := uniqueIdsMap[redId]; !exists {
				uniqueIdsMap[redId] = struct{}{}
				ids = append(ids, redId)
			}

			if _, exists := uniqueIdsMap[blueId]; !exists {
				uniqueIdsMap[blueId] = struct{}{}
				ids = append(ids, blueId)
			}
		}
	}

	return ids
}

func (c *Controller) getFightersList(fighters []*fightersmodel.Fighter) map[int32]*fightersmodel.Fighter {
	list := make(map[int32]*fightersmodel.Fighter)

	for _, fighter := range fighters {
		list[fighter.FighterId] = fighter
	}

	return list
}

func (c *Controller) eventsPretify(events []*eventmodel.Event, fighters []*fightersmodel.Fighter) []*gatewaymodel.Event {
	fightersList := c.getFightersList(fighters)
	updatedEvents := make([]*gatewaymodel.Event, len(events))

	for i, v := range events {
		event := gatewaymodel.ServiceEventToGatewayEvent(v, fightersList)
		updatedEvents[i] = event
	}

	return updatedEvents
}

func (c *Controller) GetAuthServiceHealthStatus() *model.HealthStatus {
	status, err := c.authGateway.ServiceHealthCheck()
	if err != nil {
		return badHealthStatus("auth-service")
	}

	return status
}

func (c *Controller) GetEventServiceHealthStatus() *model.HealthStatus {
	status, err := c.eventGateway.ServiceHealthCheck()
	if err != nil {
		return badHealthStatus("event-service")
	}

	return status
}

func (c *Controller) GetFightersServiceHealthStatus() *model.HealthStatus {
	status, err := c.fightersGateway.ServiceHealthCheck()
	if err != nil {
		return badHealthStatus("fighters-service")
	}

	return status
}

func badHealthStatus(serviceName string) *model.HealthStatus {
	return &model.HealthStatus{
		AppDevVersion: version.DevVersion,
		AppName:       serviceName,
		Timestamp:     time.Now().Format(time.RFC1123),
		AppRunDate:    0,
		AppTimeAlive:  0,
		Healthy:       false,
		Message:       fmt.Sprintf("[%s]: I can't feel my legs!", serviceName),
	}
}
