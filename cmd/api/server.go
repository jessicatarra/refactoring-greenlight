package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"net/http"
	"sync"
	"time"
)

type module struct {
	server *http.Server
	logger *jsonlog.Logger
}

func (m module) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		m.logger.PrintInfo("Starting Module1 server", map[string]string{"module": "legacy", "addr": m.server.Addr})
		err := m.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			m.logger.PrintFatal(err, map[string]string{
				"error": "legacy module encountered an error",
			})
		}
		m.logger.PrintInfo("Stopped Module server", map[string]string{"module": "legacy", "addr": m.server.Addr})

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

func NewModule(cfg config.Config, routes http.Handler, logger *jsonlog.Logger) *module {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      routes,
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &module{server: srv, logger: logger}
}
