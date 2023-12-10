package service

import (
	"context"
	"database/sql"
	"github.com/jessicatarra/greenlight/internal/config"
	appl "github.com/jessicatarra/greenlight/ms/auth/application"
	repo "github.com/jessicatarra/greenlight/ms/auth/repositories"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"sync"
)

func Routes(db *sql.DB, cfg config.Config, ctx context.Context, wg *sync.WaitGroup) http.Handler {

	router := httprouter.New()

	RegisterHandlers(appl.NewAppl(repo.NewUserRepo(db), repo.NewTokenRepo(db), repo.NewPermissionRepo(db), ctx, wg, cfg), router)

	return router
}
