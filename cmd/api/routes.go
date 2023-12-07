package main

import (
	"database/sql"
	"expvar"
	_ "github.com/jessicatarra/greenlight/docs"
	_authApp "github.com/jessicatarra/greenlight/ms/auth/application"
	_authRepo "github.com/jessicatarra/greenlight/ms/auth/repositories"
	_authService "github.com/jessicatarra/greenlight/ms/auth/service"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *application) routes(db *sql.DB) http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

	_authService.RegisterHandlers(_authApp.NewAppl(_authRepo.NewUserRepo(db), _authRepo.NewTokenRepo(db), app.config), router)

	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	router.Handler(http.MethodGet, "/swagger/:any", httpSwagger.WrapHandler)

	return alice.New(app.recoverPanic, app.rateLimit, app.logRequest, app.authenticate, app.enableCORS).Then(router)
}
