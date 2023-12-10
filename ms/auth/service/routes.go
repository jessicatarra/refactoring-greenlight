package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/config"
	appl "github.com/jessicatarra/greenlight/ms/auth/application"
	repo "github.com/jessicatarra/greenlight/ms/auth/repositories"
	"github.com/julienschmidt/httprouter"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	defaultIdleTimeout  = time.Minute
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

type module struct {
	server *http.Server
	logger *slog.Logger
}

func (m module) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.logger.Info("Starting Module1 server", slog.Group("server", "addr", m.server.Addr))

		err := m.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			m.logger.Info("module 1 encountered an error")

			os.Exit(1)
		}
		m.logger.Info("Stopped auth Module server", m.server.Addr)

	}()
}

func (m module) Shutdown(ctx context.Context, cancel func()) {
	defer cancel()

	err := m.server.Shutdown(ctx)
	if err != nil {
		return
	}
}

func NewModule(db *sql.DB, cfg config.Config, wg *sync.WaitGroup) *module {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8082),
		Handler:      Routes(db, cfg, wg),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &module{server: srv, logger: logger}
}

func Routes(db *sql.DB, cfg config.Config, wg *sync.WaitGroup) http.Handler {

	router := httprouter.New()

	RegisterHandlers(appl.NewAppl(repo.NewUserRepo(db), repo.NewTokenRepo(db), repo.NewPermissionRepo(db), wg, cfg), router)

	return router
}
