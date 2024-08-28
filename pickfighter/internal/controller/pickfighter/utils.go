package pickfighter

import (
	eventmodel "pickfighter.com/events/pkg/model"
	gatewaymodel "pickfighter.com/pickfighter/pkg/model"
	fightersmodel "pickfighter.com/fighters/pkg/model"
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
