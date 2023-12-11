package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/config"
	_service "github.com/jessicatarra/greenlight/ms/auth/internal/service"
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
		m.logger.Info("Starting Auth Module server", slog.Group("server", "addr", m.server.Addr))

		err := m.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			m.logger.Info("Auth module encountered an error")

			os.Exit(1)
		}
		m.logger.Info("Stopped auth Module server", slog.Group("server", "addr", m.server.Addr))

	}()
}

func (m module) Shutdown(ctx context.Context, cancel func()) {
	defer cancel()

	err := m.server.Shutdown(ctx)
	if err != nil {
		return
	}
}

func NewModule(db *sql.DB, cfg config.Config, wg *sync.WaitGroup, logger *slog.Logger) *module {
	service := _service.NewService(db, cfg, wg, logger)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 8082),
		Handler:      service.Routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &module{server: srv, logger: logger}
}
