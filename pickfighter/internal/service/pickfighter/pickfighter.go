package service

import (
	"context"
	"fmt"
	"os"

	"github.com/DoRightt/pickfighter-server/pickfighter/pkg/version"
	redisRegistry "github.com/DoRightt/pickfighter-server/pkg/discovery/redis"
	logs "github.com/DoRightt/pickfighter-server/pkg/logger"
	"github.com/DoRightt/pickfighter-server/pkg/utils"
)

var ErrAuthCertsPathRequired = fmt.Errorf("authentication certificates path is required")

type HttpHandler interface {
	RunHTTPServer(ctx context.Context) error
}

type ApiService struct {
	ServiceName string
	InstanceID  string
	Registry    *redisRegistry.Registry
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
	err := s.Registry.Deregister(ctx, s.InstanceID, s.ServiceName)
	if err != nil {
		logs.Errorf("Failed to deregister service from registry: %v", err)
	}

	logs.Warnf("Graceful shutdown. Signal received: %s", sig)

	os.Exit(0)
}
