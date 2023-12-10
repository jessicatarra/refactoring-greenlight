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
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	defaultIdleTimeout  = time.Minute
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	//defaultShutdownPeriod = 30 * time.Second
)

type module struct {
	server *http.Server
	//logger slog.Logger
}

func (m module) Start(wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		//m.logger.Info("Starting Module1 server", slog.Group("server", "addr", m.server.Addr))
		fmt.Printf("starting auth module %s", m.server.Addr)

		err := m.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			//m.logger.Info("module 1 encountered an error")
			fmt.Print("auth module encountered an error")

			os.Exit(1)
		}
		fmt.Printf("Stopped auth Module server %s", m.server.Addr)

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
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", 8082),
		Handler: Routes(db, cfg, wg),
		//ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelWarn),
		IdleTimeout:  defaultIdleTimeout,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
	}

	return &module{server: srv}
}

func Routes(db *sql.DB, cfg config.Config, wg *sync.WaitGroup) http.Handler {

	router := httprouter.New()

	RegisterHandlers(appl.NewAppl(repo.NewUserRepo(db), repo.NewTokenRepo(db), repo.NewPermissionRepo(db), wg, cfg), router)

	return router
}
