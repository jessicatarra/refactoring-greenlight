package service

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/config"
	_errors "github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"net/http"
	"strings"
)

type Middlewares interface {
	authenticate(next http.Handler) http.Handler
	requireActivatedUser(next http.HandlerFunc) http.HandlerFunc
	requirePermission(code string, next http.HandlerFunc) http.HandlerFunc
}

type middlewares struct {
	cfg  *config.Config
	appl domain.Appl
}

func (m middlewares) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")

			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				token := headerParts[1]

				user, err := m.appl.ValidateAuthTokenUseCase(token)
				if err != nil {
					_errors.ServerError(w, r, err)
					return
				}
				if user != nil {
					r = contextSetAuthenticatedUser(r, user)
				}
			}
		}

		next.ServeHTTP(w, r)
	})

}

func (m middlewares) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := contextGetAuthenticatedUser(r)

		if authenticatedUser == nil {
			_errors.AuthenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m middlewares) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		user := contextGetAuthenticatedUser(request)

		err := m.appl.UserPermissionUseCase(code, user.ID)
		if err != nil {
			if errors.Is(err, domain.ErrPermissionNotIncluded) {
				_errors.NotPermitted(writer, request)
				return
			}
			_errors.ServerError(writer, request, err)
			return
		}

		next.ServeHTTP(writer, request)
	}

	return m.requireActivatedUser(fn)
}

func (s service) Middlewares(appl domain.Appl, cfg *config.Config) Middlewares {
	return &middlewares{
		appl: appl,
		cfg:  cfg,
	}
}
