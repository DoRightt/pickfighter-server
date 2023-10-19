package services

import (
	"context"
	"os"
	"projects/fb-server/pkg/cfg"
	"projects/fb-server/pkg/pgxs"
	"sync"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ApiHandler struct {
	ServiceName string
	Router      *mux.Router
	Logger      *zap.SugaredLogger
	Repo        *pgxs.Repo
}

func New(lg *zap.SugaredLogger, name string) *ApiHandler {
	h := &ApiHandler{
		ServiceName: name,
		Logger:      lg,
		Router:      mux.NewRouter(),
	}

	return h
}

func (h *ApiHandler) Init(ctx context.Context) error {
	db, err := pgxs.NewPool(ctx, h.Logger, cfg.ViperPostgres())
	if err != nil {
		h.Logger.Errorf("Unable to start postgresql connection: %s", err)
		return err
	}
	h.Repo = db

	return nil
}

func (h *ApiHandler) Run(ctx context.Context) error {
	return h.RunHTTPServer(ctx)
}

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
