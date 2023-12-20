package model

type Fight struct {
	FightId       int32 `json:"fight_id"`
	EventId       int32 `json:"event_id"`
	FighterRedId  int32 `json:"fighter_red_id"`
	FighterBlueId int32 `json:"fighter_blue_id"`
	IsDone        bool  `json:"is_done"`
	IsCanceled    bool  `json:"is_canceled"`
	NotContest    bool  `json:"not_contest"`
	Result        int32 `json:"result"`
	CreatedAt     int64 `json:"created_at"`
	FightDate     int   `json:"fight_date"`
}

type FightResponse struct {
	FightId     int32      `json:"fight_id"`
	EventId     int32      `json:"event_id,omitempty"`
	IsDone      bool       `json:"is_done"`
	IsCanceled  bool       `json:"is_canceled"`
	NotContest  bool       `json:"not_contest"`
	Result      int32      `json:"result"`
	CreatedAt   int64      `json:"created_at"`
	FightDate   int        `json:"fight_date,omitempty"`
	FighterRed  FighterReq `json:"fighter_red"`
	FighterBlue FighterReq `json:"fighter_blue"`
}
