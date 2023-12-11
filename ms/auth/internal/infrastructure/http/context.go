package http

import (
	"context"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"net/http"
)

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

func contextSetAuthenticatedUser(r *http.Request, user *domain.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func contextGetAuthenticatedUser(r *http.Request) *domain.User {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*domain.User)
	if !ok {
		return nil
	}

	return user
}
