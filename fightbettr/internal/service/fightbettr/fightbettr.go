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

func New(h HttpHandler) ApiService {
	logger := lg.GetSugared()

	return ApiService{
		ServiceName: version.Name,
		Handler:     h,
		Logger:      logger,
	}
}

func (s *ApiService) Run(ctx context.Context) error {
	return s.Handler.RunHTTPServer(ctx)
}

func (s *ApiService) GracefulShutdown(ctx context.Context, sig string) {
	s.Logger.Warnf("Graceful shutdown. Signal received: %s", sig)

	os.Exit(0)
}
