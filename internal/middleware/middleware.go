package middleware

import (
	"fmt"
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/errors"
	"golang.org/x/time/rate"
	"net/http"
)

type Middleware interface {
	RecoverPanic(next http.Handler) http.Handler
	RateLimit(next http.Handler) http.Handler
	EnableCORS(next http.Handler) http.Handler
}

type middleware struct {
	cfg *config.Config
}

func NewSharedMiddleware(cfg *config.Config) Middleware {
	return &middleware{cfg: cfg}
}

func (m *middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {

				writer.Header().Set("Connection", "close")

				errors.ServerError(writer, request, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func (m *middleware) RateLimit(next http.Handler) http.Handler {
	//limiter := rate.NewLimiter(10, 40)
	limiter := rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !limiter.Allow() {
			errors.RateLimitExceeded(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (m *middleware) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Origin")

		writer.Header().Add("Vary", "Access-Control-Request-Method")

		origin := request.Header.Get("Origin")

		if origin != "" && len(m.cfg.Cors.TrustedOrigins) != 0 {
			for i := range m.cfg.Cors.TrustedOrigins {
				if origin == m.cfg.Cors.TrustedOrigins[i] {
					writer.Header().Set("Access-Control-Allow-Origin", origin)

					if request.Method == http.MethodOptions && request.Header.Get("Access-Control-Request-Method") != "" {

						writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
						writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

						writer.WriteHeader(http.StatusOK)
						return
					}
				}
			}
		}

		next.ServeHTTP(writer, request)
	})
}
