package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/config"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

type module struct {
	server *http.Server
	logger *slog.Logger
}

func (m module) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.logger.Info("Starting Legacy Module server", slog.Group("server", "addr", m.server.Addr))
		err := m.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			m.logger.Info("legacy module encountered an error")

			os.Exit(1)
		}
		m.logger.Info("Stopped legacy Module server", slog.Group("server", "addr", m.server.Addr))

	}()
}

func (m module) Shutdown(ctx context.Context, cancel func()) {
	defer cancel()

	err := m.server.Shutdown(ctx)
	if err != nil {
		return
	}
}

const (
	defaultIdleTimeout  = time.Minute
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
)

func NewModule(cfg config.Config, routes http.Handler, logger *slog.Logger) *module {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      routes,
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &module{server: srv, logger: logger}
}
