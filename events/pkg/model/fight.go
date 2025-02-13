package model

import (
	"database/sql"
	"time"

	fightersmodel "github.com/DoRightt/pickfighter-server/fighters/pkg/model"
)

// Fight is a structure with fight information and fighters ids
type Fight struct {
	FightId       int32        `json:"fight_id"`
	EventId       int32        `json:"event_id"`
	FighterRedId  int32        `json:"fighter_red_id"`
	FighterBlueId int32        `json:"fighter_blue_id"`
	IsDone        bool         `json:"is_done"`
	IsDraw        bool         `json:"is_draw"`
	IsCanceled    bool         `json:"is_canceled"`
	NotContest    bool         `json:"not_contest"`
	WinnerId      int32        `json:"winner_id,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
	FightDate     sql.NullTime `json:"fight_date"`
}

// Fight is a structure with information about the fight and contains the structures of the participating fighters
type FightResponse struct {
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
