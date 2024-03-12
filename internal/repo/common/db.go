package repo

import (
	"context"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"

	"github.com/jackc/pgx/v5"
)

const sep = ` AND `

// FbCommonRepo an interface for interacting with fight-related data in the database.
type FbCommonRepo interface {
	pgxs.FbRepo
	TxCreateEvent(ctx context.Context, tx pgx.Tx, e *model.EventsRequest) (int32, error)
	SearchEventsCount(ctx context.Context) (int32, error)
	SearchEvents(ctx context.Context) ([]*model.FullEventResponse, error)
	GetEventId(ctx context.Context, tx pgx.Tx, fightId int32) (int32, error)
	GetUndoneFightsCount(ctx context.Context, tx pgx.Tx, eventId int32) (int, error)
	SetEventDone(ctx context.Context, tx pgx.Tx, eventId int32) error

	SearchFightersCount(ctx context.Context, req *model.FightersRequest) (int32, error)
	SearchFighters(ctx context.Context, req *model.FightersRequest) ([]*model.Fighter, error)
	// performFightersQuery(req *model.FightersRequest) []string

	SearchBetsCount(ctx context.Context, userId int32) (int32, error)
	SearchBets(ctx context.Context, userId int32) ([]*model.Bet, error)
	CreateBet(ctx context.Context, bet *model.Bet) (int32, error)

	TxCreateEventFight(ctx context.Context, tx pgx.Tx, f model.Fight) error
	SetFightResult(ctx context.Context, tx pgx.Tx, fr *model.FightResultRequest) error
}

// CommonRepo represents a repository for for interacting with fight-related data in the database.
// It embeds the pgxs.Repo, which provides the basic PostgreSQL database operations.
type CommonRepo struct {
	pgxs.FbRepo
}

// New creates and returns a new instance of CommonRepo using the provided pgxs.Repo.
func New(r pgxs.FbRepo) *CommonRepo {
	return &CommonRepo{
		FbRepo: r,
	}
}
