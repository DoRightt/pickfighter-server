package model

import (
	fightersmodel "pickfighter.com/fighters/pkg/model"
)

// EventResponse represents a event response with []Event
type EventsResponse struct {
	Count  int32    `json:"count"`
	Events []*Event `json:"events"`
}

// Event represents a event struct with []Fights
type Event struct {
	EventId int32   `json:"event_id"`
	Name    string  `json:"name"`
	Fights  []Fight `json:"fights"`
	IsDone  bool    `json:"is_done"`
}

// Fight is a structure with information about the fight and contains the structures of the participating fighters
type Fight struct {
	FightId     int32                 `json:"fight_id"`
	EventId     int32                 `json:"event_id,omitempty"`
	FighterRed  fightersmodel.Fighter `json:"fighter_red"`
	FighterBlue fightersmodel.Fighter `json:"fighter_blue"`
	IsDone      bool                  `json:"is_done"`
	IsCanceled  bool                  `json:"is_canceled"`
	NotContest  bool                  `json:"not_contest"`
	Result      int32                 `json:"result"`
	CreatedAt   int64                 `json:"created_at"`
	FightDate   int                   `json:"fight_date,omitempty"`
}
