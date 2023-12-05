package main

import (
	"expvar"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/jessicatarra/greenlight/internal/jsonlog"
	"github.com/jessicatarra/greenlight/internal/mailer"
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

	app := &application{
		config: cfg,
		logger: logger,
		models: database.NewModels(db),
		mailer: mailer.New(cfg.Smtp.Host, cfg.Smtp.Port, cfg.Smtp.Username, cfg.Smtp.Password, cfg.Smtp.Sender),
	}

	err = app.serve(db)
	if err != nil {
		logger.PrintFatal(err, nil)
		return
	}
}
