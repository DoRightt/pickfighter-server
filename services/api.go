package services

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type ApiHandler struct {
	Router *mux.Router
	Logger *zap.Logger
	
}