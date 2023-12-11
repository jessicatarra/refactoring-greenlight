package service

import (
	"github.com/jessicatarra/greenlight/internal/middleware"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (s service) Routes() http.Handler {
	router := httprouter.New()

	s.Handlers(router)

	m := middleware.NewSharedMiddleware(&s.cfg, s.logger)

	return m.RecoverPanic(m.RateLimit(m.EnableCORS(m.LogRequest(router))))
}
