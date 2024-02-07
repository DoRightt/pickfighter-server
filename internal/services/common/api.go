package common

import (
	"context"
	"net/http"
	authRepo "projects/fb-server/internal/repo/auth"
	commonRepo "projects/fb-server/internal/repo/common"
	"projects/fb-server/internal/services"
)

type service struct {
	*services.ApiHandler
	Repo     *commonRepo.CommonRepo `json:"-" yaml:"-"`
	AuthRepo *authRepo.AuthRepo     `json:"-" yaml:"-"`
}

func New(h *services.ApiHandler) services.ApiService {
	return &service{
		ApiHandler: h,
		Repo:       commonRepo.New(h.Repo),
		AuthRepo:   authRepo.New(h.Repo),
	}
}

func (s *service) Shutdown(ctx context.Context, sig string) {}

func (s *service) Init(ctx context.Context) error {
	return nil
}

func (s *service) ApplyRoutes() {
	s.Router.HandleFunc("/fighters", s.CheckIsAdmin(s.SearchFighters)).Methods(http.MethodGet)

	s.Router.HandleFunc("/create/event", s.CheckIsAdmin(s.HandleNewEvent)).Methods(http.MethodPost)
	s.Router.HandleFunc("/events", s.GetEvents).Methods(http.MethodGet)

	s.Router.HandleFunc("/create/bet", s.IfLoggedIn(s.CreateBet)).Methods(http.MethodPost)
	s.Router.HandleFunc("/bets", s.IfLoggedIn(s.GetBets)).Methods(http.MethodGet)

	s.Router.HandleFunc("/create/result", s.CheckIsAdmin(s.AddResult)).Methods(http.MethodPost)
}
