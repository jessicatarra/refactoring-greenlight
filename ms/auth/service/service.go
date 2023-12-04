package service

import (
	"github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/internal/request"
	"github.com/jessicatarra/greenlight/internal/response"
	"github.com/jessicatarra/greenlight/ms/auth/app"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type resource struct {
	app app.App
}

func RegisterHandlers(app app.App, router *httprouter.Router) {
	res := &resource{app}

	router.HandlerFunc(http.MethodPost, "/v1/users", res.create)
}

// @Summary Register User
// @Description Registers a new user.
// @Tags Users
// @Accept json
// @Produce  json
// @Param name body app.CreateUserRequest true "User registration data"
// @Success 201 {object} entity.User
// @Router /users [post]
func (r *resource) create(res http.ResponseWriter, req *http.Request) {
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

	err = response.JSON(res, http.StatusCreated, map[string]*entity.User{"user": user})
	if err != nil {
		errors.ServerError(res, req, err)
	}
}
