package http

import (
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
)

type Service interface {
	Routes() http.Handler
	Handlers(router *httprouter.Router)
	Middlewares() Middlewares
}

type service struct {
	appl   domain.Appl
	cfg    config.Config
	logger *slog.Logger
}

func NewService(appl domain.Appl, cfg config.Config, logger *slog.Logger) Service {
	return &service{
		appl:   appl,
		cfg:    cfg,
		logger: logger,
	}
}
