package service

import (
	"database/sql"
	"github.com/jessicatarra/greenlight/internal/config"
	appl "github.com/jessicatarra/greenlight/ms/auth/internal/application"
	repo "github.com/jessicatarra/greenlight/ms/auth/internal/repositories"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
	"sync"
)

type Service interface {
	Routes() http.Handler
	logRequestMiddleware(next http.Handler) http.Handler
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

func (s service) Routes() http.Handler {
	router := httprouter.New()

	RegisterHandlers(appl.NewAppl(repo.NewUserRepo(s.db), repo.NewTokenRepo(s.db), repo.NewPermissionRepo(s.db), s.wg, s.cfg), router)

	return s.logRequestMiddleware(router)
}
