package model

// EventsRequest represents a request for events with name and slice of fights
type EventsRequest struct {
	Name   string `json:"name"`
	Fights []Fight
}

// EventResponse represents a event response with []Fights
type EventResponse struct {
	EventId int32   `json:"event_id"`
	Name    string  `json:"name"`
	Fights  []Fight `json:"fights"`
	IsDone  bool    `json:"is_done"`
}

// EventResponse represents a event response with []FightsResponse
type FullEventResponse struct {
	EventId int32           `json:"event_id"`
	Name    string          `json:"name"`
	Fights  []FightResponse `json:"fights"`
	IsDone  bool            `json:"is_done"`
}

// FightResultRequest represents a request for fight result with fight id, winner id and not contest flag.
type FightResultRequest struct {
	FightId    int32 `json:"fight_id"`
	WinnerId   int32 `json:"winner_id"`
	NotContest bool  `json:"not_contest"`
}
