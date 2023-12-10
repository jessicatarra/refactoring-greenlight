package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jessicatarra/greenlight/internal/config"
	"net/http"
	"os"
	"sync"
	"time"
)

type module struct {
	server *http.Server
}

func (m module) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		//m.logger.PrintInfo("Starting Module1 server", map[string]string{"module": "legacy", "addr": m.server.Addr})
		fmt.Printf("starting legacy module %s", m.server.Addr)
		err := m.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			//m.logger.PrintInfo("legacy module encountered an error", nil)
			fmt.Print("legacy module encountered an error")
			os.Exit(1)
		}
		//m.logger.PrintInfo("Stopped Module server", map[string]string{"module": "legacy", "addr": m.server.Addr})

		fmt.Printf("Stopped Module server %s", m.server.Addr)
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
	//defaultShutdownPeriod = 30 * time.Second
)

func NewModule(cfg config.Config, routes http.Handler) *module {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: routes,
		//ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &module{server: srv}
}
