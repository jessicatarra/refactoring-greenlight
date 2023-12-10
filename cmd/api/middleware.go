package main

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/pascaldekloe/jwt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		app.logger.PrintInfo("request", map[string]string{
			"request_remote_addr":     request.RemoteAddr,
			"request_proto":           request.Proto,
			"request_method":          request.Method,
			"request_url_request_uri": request.URL.RequestURI(),
		})

		next.ServeHTTP(writer, request)
	})
}

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
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}
		token := headerParts[1]

		claims, err := jwt.HMACCheck([]byte(token), []byte(app.config.Jwt.Secret))
		if err != nil {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		if !claims.Valid(time.Now()) {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}
		if claims.Issuer != "greenlight.tarralva.com" {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}
		if !claims.AcceptAudience("greenlight.tarralva.com") {
			app.invalidAuthenticationTokenResponse(writer, request)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.serverErrorResponse(writer, request, err)
			return
		}
		user, err := app.models.Users.Get(userID)

		if err != nil {
			switch {
			case errors.Is(err, database.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(writer, request)

			default:
				app.serverErrorResponse(writer, request, err)
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
			app.authenticationRequiredResponse(writer, request)
			return
		}

		if !user.Activated {
			app.inactiveAccountResponse(writer, request)
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
			app.serverErrorResponse(writer, request, err)
			return
		}
		if !permissions.Include(code) {
			app.notPermittedResponse(writer, request)
			return
		}

		next.ServeHTTP(writer, request)
	}

	return app.requireActivatedUser(fn)
}
