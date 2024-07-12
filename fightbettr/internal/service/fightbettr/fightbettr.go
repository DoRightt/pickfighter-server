package service

import (
	"context"
	"fmt"
	"os"

	"fightbettr.com/fightbettr/pkg/version"
	logs "fightbettr.com/pkg/logger"
	"fightbettr.com/pkg/utils"
)

var ErrAuthCertsPathRequired = fmt.Errorf("authentication certificates path is required")

type HttpHandler interface {
	RunHTTPServer(ctx context.Context) error
}

type ApiService struct {
	ServiceName string
	Handler     HttpHandler
}

// New gets logger and returns new instance of ApiService
func New(h HttpHandler) ApiService {
	return ApiService{
		ServiceName: version.Name,
		Handler:     h,
	}
}

// Run starts the API service's HTTP server.
func (s *ApiService) Run(ctx context.Context) error {
	if err := utils.LoadJwtCerts(); err != nil {
		logs.Errorf("Unable to load JWT certificates: %s", err)
		return err
	}

	return s.Handler.RunHTTPServer(ctx)
}

// GracefulShutdown logs the received signal and exits the service.
func (s *ApiService) GracefulShutdown(ctx context.Context, sig string) {
	logs.Warnf("Graceful shutdown. Signal received: %s", sig)

	os.Exit(0)
}
