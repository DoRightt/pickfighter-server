package services

import (
	"context"
	"fmt"
	"net/http"
	"projects/fb-server/pkg/httplib"
	"strings"

	"github.com/spf13/viper"
)

var allowedHeaders = []string{
	"Accept",
	"Content-Type",
	"Content-Length",
	"Authorization",
	"X-Requested-With",
	// "X-HTTP-Method-Override",
	"Cookie",
}

var allowedMethods = []string{
	http.MethodOptions,
	http.MethodGet,
	http.MethodPost,
	http.MethodPatch,
	http.MethodPut,
	http.MethodDelete,
}

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

	h.Logger.Infow("Handling request", "method", r.Method, "path", r.URL.Path, "query", r.URL.RawQuery)

	h.Router.ServeHTTP(w, r.WithContext(ctx))
}

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
	if len(srvAddr) < 1 || strings.Index(srvAddr, ":") < 0 {
		return fmt.Errorf("'%s' service address not specified", h.ServiceName)
	}

	h.Logger.Infof("Start listen '%s' http: %s", h.ServiceName, srvAddr)
	fmt.Printf("Server is listening at: %s\n", srvAddr)
	return http.ListenAndServe(srvAddr, h)
}
