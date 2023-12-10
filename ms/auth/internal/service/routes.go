package service

import (
	"database/sql"
	"github.com/jessicatarra/greenlight/internal/config"
	appl "github.com/jessicatarra/greenlight/ms/auth/internal/application"
	repo "github.com/jessicatarra/greenlight/ms/auth/internal/repositories"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

func Routes(db *sql.DB, cfg config.Config, wg *sync.WaitGroup) http.Handler {

	router := httprouter.New()

	RegisterHandlers(appl.NewAppl(repo.NewUserRepo(db), repo.NewTokenRepo(db), repo.NewPermissionRepo(db), wg, cfg), router)

	return router
}
