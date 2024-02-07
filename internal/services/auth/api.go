package auth

import (
	"context"
	"net/http"
	"projects/fb-server/pkg/model"
	authRepo "projects/fb-server/internal/repo/auth"
	"projects/fb-server/internal/services"
)

type service struct {
	*services.ApiHandler

	Repo *authRepo.AuthRepo `json:"-" yaml:"-"`
}

func New(h *services.ApiHandler) services.ApiService {
	return service{
		ApiHandler: h,
		Repo:       authRepo.New(h.Repo),
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
	s.Router.HandleFunc("/register/confirm", s.ConfirmRegistration).Methods(http.MethodPost)
	s.Router.HandleFunc("/login", s.Login).Methods(http.MethodPost)
	s.Router.HandleFunc("/logout", s.IfLoggedIn(s.Logout)).Methods(http.MethodGet)
	s.Router.HandleFunc("/password/reset", s.ResetPassword).Methods(http.MethodPost)
	s.Router.HandleFunc("/password/recover", s.RecoverPassword).Methods(http.MethodPost)

	// profile
	s.Router.Handle("/profile", s.IfLoggedIn(s.GetCurrentUser)).Methods(http.MethodGet)
}
