package common

import (
	"context"
	"net/http"

	// authRepo "fightbettr.com/fb-server/internal/repo/auth"
	commonRepo "fightbettr.com/fb-server/internal/repo/common"
	"fightbettr.com/fb-server/internal/services"

	"github.com/jackc/pgx/v5"
)

type CommonService interface {
	services.ApiService
	CheckEventIsDone(ctx context.Context, tx pgx.Tx, fightId int32) error
	HandleNewEvent(w http.ResponseWriter, r *http.Request)
}

type service struct {
	*services.ApiHandler

	Repo commonRepo.FbCommonRepo `json:"-" yaml:"-"`
	// AuthRepo *authRepo.AuthRepo     `json:"-" yaml:"-"` // TODO
}

// New creates a new instance of the service using the provided ApiHandler and initializes an commonRepo for working with the common repository.
func New(h *services.ApiHandler) CommonService {
	return &service{
		ApiHandler: h,
		Repo:       commonRepo.New(h.Repo),
		// AuthRepo:   authRepo.New(h.Repo), // TODO
	}
}

func (s *service) Shutdown(ctx context.Context, sig string) {}

func (s *service) Init(ctx context.Context) error {
	return nil
}

// ApplyRoutes sets up and assigns the handlers for various API endpoints. It uses the Gorilla Mux
// router to define the routes and associate them with the corresponding handler functions. The
// routes include functionalities such as searching for fighters, creating events, retrieving events,
// creating bets, retrieving bets, and adding fight results. Access to certain routes is restricted based
// on user roles, such as admin or logged-in user. The CheckIsAdmin and IfLoggedIn middleware functions
// are used to enforce role-based access control for specific routes.
func (s *service) ApplyRoutes() {
	s.Router.HandleFunc("/fighters", s.CheckIsAdmin(s.SearchFighters)).Methods(http.MethodGet)

	s.Router.HandleFunc("/create/event", s.CheckIsAdmin(s.HandleNewEvent)).Methods(http.MethodPost)
	s.Router.HandleFunc("/events", s.GetEvents).Methods(http.MethodGet)

	s.Router.HandleFunc("/create/bet", s.IfLoggedIn(s.CreateBet)).Methods(http.MethodPost)
	s.Router.HandleFunc("/bets", s.IfLoggedIn(s.GetBets)).Methods(http.MethodGet)

	s.Router.HandleFunc("/create/result", s.CheckIsAdmin(s.AddResult)).Methods(http.MethodPost)
}
