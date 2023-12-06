package service

import (
	"github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/internal/request"
	"github.com/jessicatarra/greenlight/internal/response"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Resource interface {
	create(res http.ResponseWriter, req *http.Request)
	activate(res http.ResponseWriter, req *http.Request)
}

type resource struct {
	appl domain.Appl
}

func RegisterHandlers(appl domain.Appl, router *httprouter.Router) {
	res := &resource{appl}

	router.HandlerFunc(http.MethodPost, "/v1/users", res.create)
}

// @Summary Register User
// @Description Registers a new user.
// @Tags Users
// @Accept json
// @Produce  json
// @Param name body domain.CreateUserRequest true "User registration data"
// @Success 201 {object} domain.User
// @Router /users [post]
func (r *resource) create(res http.ResponseWriter, req *http.Request) {
	var input domain.CreateUserRequest

	err := request.DecodeJSON(res, req, &input)
	if err != nil {
		errors.BadRequest(res, req, err)
		return
	}

	existingUser, _ := r.appl.GetByEmailUseCase(input)

	ValidateUser(input, existingUser)

	if input.Validator.HasErrors() {
		errors.FailedValidation(res, req, input.Validator)
		return
	}

	user, err := r.appl.CreateUseCase(input)
	if err != nil {
		errors.ServerError(res, req, err)
		return
	}

	err = response.JSON(res, http.StatusCreated, map[string]*domain.User{"user": user})
	if err != nil {
		errors.ServerError(res, req, err)
	}
}

func (r *resource) activate(res http.ResponseWriter, req *http.Request) {
	//TODO implement me
	panic("implement me")
}
