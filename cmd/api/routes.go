package main

import (
	"expvar"
	_ "github.com/jessicatarra/greenlight/docs"
	"github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/internal/middleware"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(errors.NotFound)

	router.MethodNotAllowed = http.HandlerFunc(errors.MethodNotAllowed)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission("movies:read", app.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission("movies:write", app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.requirePermission("movies:read", app.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", app.requirePermission("movies:write", app.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.requirePermission("movies:write", app.deleteMovieHandler))

	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	router.Handler(http.MethodGet, "/swagger/:any", httpSwagger.WrapHandler)

	m := middleware.NewSharedMiddleware(&app.config, app.logger)

	return m.RecoverPanic(m.RateLimit(m.EnableCORS(app.authenticate(m.LogRequest(router)))))
}
