package model

import (
	"database/sql"
	"time"

	fightersmodel "github.com/DoRightt/pickfighter-server/fighters/pkg/model"
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
	IsDraw      bool                  `json:"is_draw"`
	IsCanceled  bool                  `json:"is_canceled"`
	NotContest  bool                  `json:"not_contest"`
	Winner_id   int32                 `json:"winner_id,omitempty"`
	CreatedAt   time.Time             `json:"created_at"`
	FightDate   sql.NullTime          `json:"fight_date,omitempty"`
}
