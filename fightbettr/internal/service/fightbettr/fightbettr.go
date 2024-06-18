package service

import (
	"context"
	"os"

	lg "fightbettr.com/fightbettr/pkg/logger"
	"fightbettr.com/fightbettr/pkg/version"
)

type HttpHandler interface {
	RunHTTPServer(ctx context.Context) error
}

type ApiService struct {
	ServiceName string
	Handler     HttpHandler
	Logger      lg.FbLogger
}

// New gets logger and returns new instance of ApiService
func New(h HttpHandler) ApiService {
	logger := lg.GetSugared()

	return ApiService{
		ServiceName: version.Name,
		Handler:     h,
		Logger:      logger,
	}
}

// Run starts the API service's HTTP server.
func (s *ApiService) Run(ctx context.Context) error {
	return s.Handler.RunHTTPServer(ctx)
}

// GracefulShutdown logs the received signal and exits the service.
func (s *ApiService) GracefulShutdown(ctx context.Context, sig string) {
	s.Logger.Warnf("Graceful shutdown. Signal received: %s", sig)

	os.Exit(0)
}
