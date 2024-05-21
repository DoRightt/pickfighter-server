package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"fightbettr.com/fightbettr/internal/controller/fightbettr"
	lg "fightbettr.com/fightbettr/pkg/logger"
	"fightbettr.com/fightbettr/pkg/version"
	"fightbettr.com/pkg/httplib"
	"fightbettr.com/pkg/ipaddr"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

// allowedHeaders defines the list of allowed HTTP headers that can be used in CORS requests.
var allowedHeaders = []string{
	"Accept",
	"Content-Type",
	"Content-Length",
	"Authorization",
	"X-Requested-With",
	// "X-HTTP-Method-Override",
	"Cookie",
}

// allowedMethods defines the list of allowed HTTP methods that can be used in CORS requests.
var allowedMethods = []string{
	http.MethodOptions,
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodPut,
	http.MethodDelete,
}

// ContextKey represents a custom type for identifying context keys in HTTP requests.
type ContextKey string

// Constants representing various context keys for use in HTTP requests.
const (
	ContextKeyHost           ContextKey = "host"
	ContextKeyPath           ContextKey = "path"
	ContextKeyRemoteAddr     ContextKey = "remote_addr"
	ContextKeyCFConnectingIP ContextKey = "cf_connecting_ip"
)

// Handler defines a movie handler.
type Handler struct {
	ctrl   *fightbettr.Controller
	router *mux.Router
	logger lg.FbLogger
}

func New(ctrl *fightbettr.Controller) *Handler {
	return &Handler{
		ctrl:   ctrl,
		logger: lg.GetSugared(),
		router: mux.NewRouter(),
	}
}

// ServeHTTP handles the incoming HTTP request by setting CORS headers, processing preflight OPTIONS requests,
// and forwarding the request to the underlying router with additional context values.
// It checks for the "Origin" header to set CORS headers and responds to OPTIONS requests appropriately.
// The function logs details of the incoming request and forwards it to the router for further handling.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ","))
	}

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx = context.WithValue(ctx, ContextKeyHost, r.Host)
	ctx = context.WithValue(ctx, ContextKeyPath, r.URL.Path)
	ctx = context.WithValue(ctx, ContextKeyRemoteAddr, r.RemoteAddr)
	ctx = context.WithValue(ctx, ContextKeyCFConnectingIP, r.Header.Get(ipaddr.CFConnectingIp))

	h.logger.Infow("Handling request", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery)

	h.router.ServeHTTP(w, r.WithContext(ctx))
}

// RunHTTPServer starts the HTTP server for the API handler with specified routes and services.
// It sets the health check route and adds routes for each registered service.
// The server listens on the specified address, and if successful, it prints the server address.
// It returns an error if the service address is not specified or if there is an issue starting the server.
func (h *Handler) RunHTTPServer(ctx context.Context) error {
	serviceName := version.Name
	httplib.SetCookieName(viper.GetString("auth.cookie_name"))

	// sys routes
	// h.router.HandleFunc("/health", h.HealthCheck).Methods(http.MethodGet)

	h.ApplyRoutes()

	srvAddr := viper.GetString("http.addr")
	if len(srvAddr) < 1 || !strings.Contains(srvAddr, ":") {
		return fmt.Errorf("'%s' service address not specified", serviceName)
	}

	h.logger.Infof("Start listen '%s' http: %s", serviceName, srvAddr)
	fmt.Printf("Server is listening at: %s\n", srvAddr)

	return http.ListenAndServe(srvAddr, h)
}

// ApplyRoutes sets up the API routes for related services.
// It associates each route with the corresponding handler method from the service.
// The routes include user registration, login, logout, password reset, password recovery, and profile retrieval.
func (h *Handler) ApplyRoutes() {

	// fighters
	h.router.HandleFunc("/fighters", h.GetFighters).Methods(http.MethodGet)
}
