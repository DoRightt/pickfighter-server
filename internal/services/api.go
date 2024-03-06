package services

import (
	"context"
	"net/http"
	"os"
	"projects/fb-server/pkg/logger"
	"projects/fb-server/pkg/model"
	"projects/fb-server/pkg/pgxs"
	"sync"

	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

// FbRouter represents an interface for routing and handling HTTP requests.
type FbRouter interface {
	ServeHTTP(w http.ResponseWriter, req *http.Request)
	Handle(path string, handler http.Handler) *mux.Route
	HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route
}

// ApiService defines the interface that any API service must implement.
type ApiService interface {
	Init(ctx context.Context) error
	ApplyRoutes()
	Shutdown(ctx context.Context, sig string)
}

// ApiHandlerInterface defines the main app interface.
type AppInterface interface {
	Init(repo pgxs.FbRepo) error
	Run(ctx context.Context) error
	GracefulShutdown(ctx context.Context, sig string)
	AddService(name string, srv ApiService)
	RunHTTPServer(ctx context.Context) error
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	loadJwtCerts() error
	verifyJWT(jwtRawValue string) (jwt.Token, error)
	IfLoggedIn(fn http.HandlerFunc) http.HandlerFunc
	CheckIsAdmin(next http.HandlerFunc) http.HandlerFunc
	HealthCheck(w http.ResponseWriter, r *http.Request)
	HandleEmailEvent(ctx context.Context, data *model.EmailData)
}

// ApiHandler represents the main handler for the API. It holds information about router, logger, repository, and services.
type ApiHandler struct {
	ServiceName string
	Router      FbRouter
	Logger      logger.FbLogger
	Repo        pgxs.FbRepo

	Services map[string]ApiService `json:"-" yaml:"-"`
}

// New creates a new instance of ApiHandler with the provided logger, service name, and initializes the router and services.
func New(lg logger.FbLogger, name string) *ApiHandler {
	h := &ApiHandler{
		ServiceName: name,
		Logger:      lg,
		Router:      mux.NewRouter(),
		Services:    make(map[string]ApiService),
	}

	return h
}

// Init initializes the ApiHandler by establishing a connection to PostgreSQL using the special configuration.
// It also loads JWT certificates required for authentication.
// If any error occurs during initialization, it is logged, and the error is returned.
func (h *ApiHandler) Init(repo pgxs.FbRepo) error {
	h.Repo = repo

	if err := h.loadJwtCerts(); err != nil {
		h.Logger.Errorf("Unable to load JWT certificates: %s", err)
		return err
	}

	return nil
}

// Run initializes and starts the services registered with the ApiHandler.
// It iterates through the available services and initializes each one using the provided context.
// After initializing the services, it starts the HTTP server to handle incoming requests.
func (h *ApiHandler) Run(ctx context.Context) error {
	for name := range h.Services {
		srv, ok := h.Services[name]
		if ok {
			if err := srv.Init(ctx); err != nil {
				return err
			}
		}
	}

	return h.RunHTTPServer(ctx)
}

// AddService adds an instance of the ApiService to the ApiHandler's services map.
func (h *ApiHandler) AddService(name string, srv ApiService) {
	h.Services[name] = srv
}

// GracefulShutdown performs a graceful shutdown of the API service.
// It is triggered when a specified signal (e.g., SIGINT or SIGTERM) is received.
// The method initiates a graceful shutdown of the underlying repository (if available),
// waits for the repository's shutdown to complete, and then exits the application.
func (h *ApiHandler) GracefulShutdown(ctx context.Context, sig string) {
	h.Logger.Warnf("Graceful shutdown. Signal received: %s", sig)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		if h.Repo != nil {
			h.Repo.GracefulShutdown()
		}
	}()

	wg.Wait()

	os.Exit(0)
}
