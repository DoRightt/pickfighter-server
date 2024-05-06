package auth

import (
	"context"
	"net/http"
	authRepo "fightbettr.com/fb-server/internal/repo/auth"
	"fightbettr.com/fb-server/internal/services"
	"fightbettr.com/fb-server/pkg/model"
)

type AuthService interface {
	services.ApiService
}

type service struct {
	*services.ApiHandler

	Repo authRepo.FbAuthRepo `json:"-" yaml:"-"`
}

// New creates a new instance of the service using the provided ApiHandler and initializes an AuthRepo for working with the authentication repository.
func New(h *services.ApiHandler) AuthService {
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

// Name returns the service name
func (s service) Name() string {
	return model.AuthService
}

// ApplyRoutes sets up the API routes for the authentication and profile-related endpoints.
// It associates each route with the corresponding handler method from the service.
// The routes include user registration, login, logout, password reset, password recovery, and profile retrieval.
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
