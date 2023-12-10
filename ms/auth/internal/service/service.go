package service

import (
	"database/sql"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
	"sync"
)

type Service interface {
	Routes() http.Handler
	logRequestMiddleware(next http.Handler) http.Handler
	registerHandlers(appl domain.Appl, router *httprouter.Router)
}

type service struct {
	db     *sql.DB
	cfg    config.Config
	wg     *sync.WaitGroup
	logger *slog.Logger
}

func NewService(db *sql.DB, cfg config.Config, wg *sync.WaitGroup, logger *slog.Logger) Service {
	return &service{
		db:     db,
		cfg:    cfg,
		wg:     wg,
		logger: logger,
	}
}
