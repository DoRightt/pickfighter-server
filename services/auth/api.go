package auth

import (
	"context"
	"net/http"
	"projects/fb-server/pkg/model"
	authRepo "projects/fb-server/repo/auth"
	"projects/fb-server/services"
)

type service struct {
	*services.ApiHandler

	Repo *authRepo.AuthRepo `json:"-" yaml:"-"`
}

func New(h *services.ApiHandler) services.ApiService {
	return service{
		ApiHandler: h,
	}
}

func (s service) Init(ctx context.Context) error {
	return nil
}

func (s service) Shutdown(ctx context.Context, sig string) {

}

func (s service) Name() string {
	return model.AuthService
}

func (s service) ApplyRoutes() {
	// authentication
	s.Router.HandleFunc("/register", s.Register).Methods(http.MethodPost)
	s.Router.HandleFunc("/login", s.Login).Methods(http.MethodPost)
}
