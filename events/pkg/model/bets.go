package model

type BetsResponse struct {
	Count int32 `json:"count"`
	Bets  []*Bet `json:"bets"`
}

// Bet represents users bet properties
type Bet struct {
	BetId     int32 `json:"bet_id"`
	FightId   int32 `json:"fight_id"`
	UserId    int32 `json:"user_id"`
	FighterId int32 `json:"fighter_id"`
}
