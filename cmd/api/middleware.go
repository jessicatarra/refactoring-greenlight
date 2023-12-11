package main

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/database"
	_errors "github.com/jessicatarra/greenlight/internal/errors"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TODO: relocate following middleware into auth module and remote called then from legacy module

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Add("Vary", "Authorization")
		authorizationHeader := request.Header.Get("Authorization")
		if authorizationHeader == "" {
			request = app.contextSetUser(request, database.AnonymousUser)
			next.ServeHTTP(writer, request)
			return
		}
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			_errors.InvalidAuthenticationToken(writer, request)
			return
		}
		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.Jwt.Secret))
		if err != nil {
			_errors.InvalidAuthenticationToken(writer, request)
			return
		}

		if !claims.Valid(time.Now()) {
			_errors.InvalidAuthenticationToken(writer, request)
			return
		}
		if claims.Issuer != "greenlight.tarralva.com" {
			_errors.InvalidAuthenticationToken(writer, request)
			return
		}
		if !claims.AcceptAudience("greenlight.tarralva.com") {
			_errors.InvalidAuthenticationToken(writer, request)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			_errors.ServerError(writer, request, err)
			return
		}
		user, err := app.models.Users.Get(userID)

		if err != nil {
			switch {
			case errors.Is(err, database.ErrRecordNotFound):
				_errors.InvalidAuthenticationToken(writer, request)

			default:
				_errors.ServerError(writer, request, err)
			}
			return
		}

		request = app.contextSetUser(request, user)
		next.ServeHTTP(writer, request)
	})
}

func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := app.contextGetUser(request)

		if user.IsAnonymous() {
			_errors.AuthenticationRequired(writer, request)
			return
		}

		if !user.Activated {
			_errors.InactiveAccount(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		user := app.contextGetUser(request)

		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		if err != nil {
			_errors.ServerError(writer, request, err)
			return
		}
		if !permissions.Include(code) {
			_errors.NotPermitted(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	}

	return app.requireActivatedUser(fn)
}
