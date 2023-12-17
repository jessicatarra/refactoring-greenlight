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

func (a *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(errors.NotFound)

	router.MethodNotAllowed = http.HandlerFunc(errors.MethodNotAllowed)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", a.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/movies", a.requirePermission("movies:read", a.listMoviesHandler))
	router.HandlerFunc(http.MethodPost, "/v1/movies", a.requirePermission("movies:write", a.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", a.requirePermission("movies:read", a.showMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:id", a.requirePermission("movies:write", a.updateMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", a.requirePermission("movies:write", a.deleteMovieHandler))

	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())

	router.Handler(http.MethodGet, "/swagger/:any", httpSwagger.WrapHandler)

	m := middleware.NewSharedMiddleware(&a.config, a.logger)

	return m.RecoverPanic(m.RateLimit(m.EnableCORS(a.authenticate(m.LogRequest(router)))))
}
