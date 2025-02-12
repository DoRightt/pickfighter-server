package model

// EventRequest represents a request for event with name and slice of fights
type EventRequest struct {
	Name   string `json:"name"`
	Fights []Fight
}

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

// EventResponse represents a event response with []FightsResponse
type FullEventResponse struct {
	EventId int32           `json:"event_id"`
	Name    string          `json:"name"`
	Fights  []FightResponse `json:"fights"`
	IsDone  bool            `json:"is_done"`
}

// FightResultRequest represents a request for fight result with fight id, winner id and not contest flag.
type FightResultRequest struct {
	FightId    int32  `json:"fight_id"`
	WinnerId   int32  `json:"winner_id"`
	IsCanceled bool   `json:"is_canceled"`
	IsDraw     bool   `json:"is_draw"`
	NotContest bool   `json:"not_contest"`
	Method     string `json:"method"`
	Round      int    `json:"round"`
}
