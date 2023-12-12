package main

import (
	"context"
	pb "github.com/jessicatarra/greenlight/api/proto"
	"github.com/jessicatarra/greenlight/internal/database"
	_errors "github.com/jessicatarra/greenlight/internal/errors"
	"net/http"
	"strings"
	"time"
)

func (a *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader != "" {
			headerParts := strings.Split(authorizationHeader, " ")

			if len(headerParts) == 2 && headerParts[0] == "Bearer" {
				token := headerParts[1]

				grpcReq := &pb.ValidateAuthTokenRequest{
					Token: token,
				}
				user, err := a.grpcClient.ValidateAuthToken(context.Background(), grpcReq)

				if err != nil {
					_errors.ServerError(w, r, err)
					return
				}
				createdAt := time.Unix(user.CreatedAt.Seconds, int64(user.CreatedAt.Nanos))
				if user != nil {
					r = a.contextSetUser(r, &database.User{
						ID:        user.Id,
						CreatedAt: createdAt,
						Name:      user.Name,
						Email:     user.Email,
						//Password:  user.HashedPassword,
						Activated: user.Activated,
						Version:   int(user.Version),
					})
				}
			}
		}

		next.ServeHTTP(w, r)
	})

}

func (a *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authenticatedUser := a.contextGetUser(r)

		if authenticatedUser == nil {
			_errors.AuthenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		user := a.contextGetUser(request)

		grpcReq := &pb.UserPermissionRequest{
			Code:   code,
			UserId: user.ID,
		}
		_, err := a.grpcClient.UserPermission(context.Background(), grpcReq)
		if err != nil {
			//if errors.Is(err, domain.ErrPermissionNotIncluded) {
			//	_errors.NotPermitted(writer, request)
			//	return
			//}
			_errors.ServerError(writer, request, err)
			return
		}

		next.ServeHTTP(writer, request)
	}

	return a.requireActivatedUser(fn)
}
