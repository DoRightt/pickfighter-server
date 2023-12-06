package model

type EventsRequest struct {
	Name   string `json:"name"`
	Fights []Fight
}

type EventResponse struct {
	EventId int32  `json:"event_id"`
	Name    string `json:"name"`
	Fights  []Fight `json:"fights"`
}
