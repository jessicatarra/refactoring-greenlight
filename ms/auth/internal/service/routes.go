package service

import (
	appl "github.com/jessicatarra/greenlight/ms/auth/internal/application"
	repo "github.com/jessicatarra/greenlight/ms/auth/internal/repositories"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s service) Routes() http.Handler {
	router := httprouter.New()

	s.registerHandlers(appl.NewAppl(repo.NewUserRepo(s.db), repo.NewTokenRepo(s.db), repo.NewPermissionRepo(s.db), s.wg, s.cfg), router)

	return s.logRequestMiddleware(router)
}
