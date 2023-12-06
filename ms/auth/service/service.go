package service

import (
	"github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/internal/request"
	"github.com/jessicatarra/greenlight/internal/response"
	"github.com/jessicatarra/greenlight/ms/auth/application"
	"github.com/jessicatarra/greenlight/ms/auth/entity"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Resource interface {
	create(res http.ResponseWriter, req *http.Request)
	activate(res http.ResponseWriter, req *http.Request)
}

type resource struct {
	appl application.Appl
}

func RegisterHandlers(appl application.Appl, router *httprouter.Router) {
	res := &resource{appl}

	router.HandlerFunc(http.MethodPost, "/v1/users", res.create)
}

// @Summary Register User
// @Description Registers a new user.
// @Tags Users
// @Accept json
// @Produce  json
// @Param name body entity.CreateUserRequest true "User registration data"
// @Success 201 {object} entity.User
// @Router /users [post]
func (r *resource) create(res http.ResponseWriter, req *http.Request) {
	var input entity.CreateUserRequest

	err := request.DecodeJSON(res, req, &input)
	if err != nil {
		errors.BadRequest(res, req, err)
		return
	}

	user, err := r.appl.CreateUseCase(input)
	if err != nil {
		errors.ServerError(res, req, err)
		return
	}

	err = response.JSON(res, http.StatusCreated, map[string]*entity.User{"user": user})
	if err != nil {
		errors.ServerError(res, req, err)
	}
}

func (r *resource) activate(res http.ResponseWriter, req *http.Request) {

}
