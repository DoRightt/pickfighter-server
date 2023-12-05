package model

type EventsRequest struct {
	Name string `json:"name"`
}

type EventResponse struct {
	EventId int32  `json:"event_id"`
	Name    string `json:"name"`
}
