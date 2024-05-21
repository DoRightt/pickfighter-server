package fightbettr

import (
	"context"

	fightersmodel "fightbettr.com/fighters/pkg/model"
)

type fightersGateway interface {
	SearchFighters(ctx context.Context, status fightersmodel.FighterStatus) ([]*fightersmodel.Fighter, error)
}

// Controller defines a gateway service controller.
type Controller struct {
	fightersGateway fightersGateway
}

func New(fightersGateway fightersGateway) *Controller {
	return &Controller{
		fightersGateway,
	}
}

// SearchFighters searches for fighters with the given status using the fightersGateway.
func (c *Controller) SearchFighters(ctx context.Context, status string) ([]*fightersmodel.Fighter, error) {
	fighters, err := c.fightersGateway.SearchFighters(ctx, fightersmodel.FighterStatus(status))
	if err != nil {
		return nil, err
	}

	return fighters, nil
}
