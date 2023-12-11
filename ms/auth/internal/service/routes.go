package service

import (
	"github.com/jessicatarra/greenlight/internal/middleware"
	appl "github.com/jessicatarra/greenlight/ms/auth/internal/application"
	repo "github.com/jessicatarra/greenlight/ms/auth/internal/repositories"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s service) Routes() http.Handler {
	router := httprouter.New()

	newAppl := appl.NewAppl(repo.NewUserRepo(s.db), repo.NewTokenRepo(s.db), repo.NewPermissionRepo(s.db), s.wg, s.cfg)

	s.Handlers(newAppl, router)

	m := middleware.NewSharedMiddleware(&s.cfg, s.logger)

	s.Middlewares(newAppl, &s.cfg)

	return m.RecoverPanic(m.RateLimit(m.EnableCORS(m.LogRequest(router))))
}
