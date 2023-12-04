package service

import (
	"github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/internal/request"
	"github.com/jessicatarra/greenlight/internal/response"
	"github.com/jessicatarra/greenlight/ms/auth/app"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"net/http"
)

type resource struct {
	app app.App
}

func (r resource) Create(res http.ResponseWriter, req *http.Request) {
	var input app.CreateUserRequest

	err := request.DecodeJSON(res, req, &input)
	if err != nil {
		errors.BadRequest(res, req, err)
		return
	}

	user, err := r.app.Create(input)
	if err != nil {
		errors.ServerError(res, req, err)
	}

	err = response.JSON(res, http.StatusCreated, map[string]entity.User{"user": user})
	if err != nil {
		errors.ServerError(res, req, err)
	}
}
