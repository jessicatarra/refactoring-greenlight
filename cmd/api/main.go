package main

import (
	"database/sql"
	"expvar"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/jessicatarra/greenlight/internal/mailer"
	_auth "github.com/jessicatarra/greenlight/ms/auth"
	"log/slog"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"
)

type application struct {
	config config.Config
	logger *slog.Logger
	models database.Models
	mailer mailer.Mailer
	wg     sync.WaitGroup
}

// @title Greenlight API Docs
// @version 1.0.0
// @contact.name Jessica Tarra
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	err := run(logger)
	if err != nil {
		trace := string(debug.Stack())
		logger.Error(err.Error(), "trace", trace)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := config.Init()
	if err != nil {
		return err
	}

	db, err := database.New(cfg.DB.Dsn, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime, true)
	if err != nil {
		return err
	}
	defer db.Close()

	initMetrics(db)

	app := newLegacyApplication(cfg, logger, db)

	monolith := NewModularMonolith(&app.wg)

	monolith.AddModule(NewModule(cfg, app.routes(), app.logger))
	monolith.AddModule(_auth.NewModule(db, cfg, &app.wg, app.logger))

	return monolith.Run()
}

func initMetrics(db *sql.DB) {
	expvar.NewString("version").Set(config.Version)

	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))
}

func newLegacyApplication(cfg config.Config, logger *slog.Logger, db *sql.DB) *application {
	return &application{
		config: cfg,
		logger: logger,
		models: database.NewModels(db),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.From),
	}
}
