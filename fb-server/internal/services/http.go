package services

import (
	"context"
	"fmt"
	"net/http"
	"fightbettr.com/fb-server/pkg/httplib"
	"fightbettr.com/fb-server/pkg/ipaddr"
	"strings"

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

// ServeHTTP handles the incoming HTTP request by setting CORS headers, processing preflight OPTIONS requests,
// and forwarding the request to the underlying router with additional context values.
// It checks for the "Origin" header to set CORS headers and responds to OPTIONS requests appropriately.
// The function logs details of the incoming request and forwards it to the router for further handling.
func (h *ApiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	h.Logger.Infow("Handling request", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery)

	h.Router.ServeHTTP(w, r.WithContext(ctx))
}

// RunHTTPServer starts the HTTP server for the API handler with specified routes and services.
// It sets the health check route and adds routes for each registered service.
// The server listens on the specified address, and if successful, it prints the server address.
// It returns an error if the service address is not specified or if there is an issue starting the server.
func (h *ApiHandler) RunHTTPServer(ctx context.Context) error {
	httplib.SetCookieName(viper.GetString("auth.cookie_name"))

	// sys routes
	h.Router.HandleFunc("/health", h.HealthCheck).Methods(http.MethodGet)

	for name := range h.Services {
		srv, ok := h.Services[name]
		if ok {
			h.Logger.Infof("Adding '%s' service routes", name)
			srv.ApplyRoutes()
			// h.Services[name] = srv
		}
	}

	srvAddr := viper.GetString("http.addr")
	if len(srvAddr) < 1 || !strings.Contains(srvAddr, ":") {
		return fmt.Errorf("'%s' service address not specified", h.ServiceName)
	}

	h.Logger.Infof("Start listen '%s' http: %s", h.ServiceName, srvAddr)
	fmt.Printf("Server is listening at: %s\n", srvAddr)
	return http.ListenAndServe(srvAddr, h)
}
