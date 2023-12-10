package main

import (
	"context"
	"database/sql"
	"expvar"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"github.com/jessicatarra/greenlight/internal/mailer"
	_authRoutes "github.com/jessicatarra/greenlight/ms/auth/service"
	"os"
	"runtime"
	"sync"
	"time"
)

type application struct {
	config config.Config
	logger *jsonlog.Logger
	models database.Models
	mailer mailer.Mailer
	//wg     sync.WaitGroup
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
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	cfg, err := config.Init()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	db, err := database.New(cfg.DB.Dsn, cfg.DB.MaxOpenConns, cfg.DB.MaxIdleConns, cfg.DB.MaxIdleTime, true)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	initMetrics(db)

	var wg sync.WaitGroup

	app := newLegacyApplication(cfg, logger, db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	monolith := NewModularMonolith(ctx, &wg)

	monolith.AddModule("Legacy Module", cfg.Port, app.routes())
	monolith.AddModule("Auth Module", 8082, _authRoutes.Routes(db, cfg, ctx, &wg))

	err = monolith.Run()
	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}
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

func newLegacyApplication(cfg config.Config, logger *jsonlog.Logger, db *sql.DB) *application {
	return &application{
		config: cfg,
		logger: logger,
		models: database.NewModels(db),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.From),
	}
}
