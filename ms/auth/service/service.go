package service

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/database"
	_errors "github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/internal/request"
	"github.com/jessicatarra/greenlight/internal/response"
	"github.com/jessicatarra/greenlight/internal/utils/helpers"
	"github.com/jessicatarra/greenlight/ms/auth/domain"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type envelope map[string]interface{}

type Resource interface {
	createUser(res http.ResponseWriter, req *http.Request)
	activateUser(res http.ResponseWriter, req *http.Request)
	createToken(res http.ResponseWriter, req *http.Request)
}

type resource struct {
	appl    domain.Appl
	helpers helpers.Helpers
}

func RegisterHandlers(appl domain.Appl, router *httprouter.Router) {
	res := &resource{appl, helpers.New()}

	router.HandlerFunc(http.MethodPost, "/v1/users", res.createUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", res.activateUser)
}

// @Summary Register User
// @Description Registers a new user.
// @Tags Users
// @Accept json
// @Produce  json
// @Param name body domain.CreateUserRequest true "User registration data"
// @Success 201 {object} domain.User
// @Router /users [post]
func (r *resource) createUser(res http.ResponseWriter, req *http.Request) {
	var input domain.CreateUserRequest

	err := request.DecodeJSON(res, req, &input)
	if err != nil {
		_errors.BadRequest(res, req, err)
		return
	}

	existingUser, err := r.appl.GetByEmailUseCase(input.Email)
	if err != nil && err.Error() != domain.ErrRecordNotFound.Error() {
		_errors.ServerError(res, req, err)
		return
	}

	ValidateUser(input, existingUser)

	if input.Validator.HasErrors() {
		_errors.FailedValidation(res, req, input.Validator)
		return
	}

	user, err := r.appl.CreateUseCase(input)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrDuplicateEmail):
			input.Validator.AddError("email a user with this email address already exists")
			_errors.FailedValidation(res, req, input.Validator)
		default:
			_errors.ServerError(res, req, err)
		}
		return
	}

	err = response.JSON(res, http.StatusCreated, envelope{"user": user})
	if err != nil {
		_errors.ServerError(res, req, err)
	}
}

// @Summary Activate User
// @Description Activates a user account using a token that was previously sent when successfully register a new user
// @Tags Users
// @Accept json
// @Produce  json
// @Param token query string true "Token for user activation"
// @Success 200 {object} domain.User
// @Router /users/activated [put]
func (r *resource) activateUser(res http.ResponseWriter, req *http.Request) {
	var input domain.ActivateUserRequest

	qs := req.URL.Query()

	input.TokenPlaintext = r.helpers.ReadString(qs, "token", "")

	input.Validator.Check(input.TokenPlaintext != "", "token must be provided")
	input.Validator.Check(len(input.TokenPlaintext) == 26, "token must be 26 bytes long")

	if input.Validator.HasErrors() {
		_errors.FailedValidation(res, req, input.Validator)
		return
	}

	user, err := r.appl.ActivateUseCase(input.TokenPlaintext)
	if err != nil {
		_errors.ServerError(res, req, err)
		return
	}

	err = response.JSON(res, http.StatusCreated, envelope{"user": user})
	if err != nil {
		_errors.ServerError(res, req, err)
	}
}

func (r *resource) createToken(res http.ResponseWriter, req *http.Request) {
	var input domain.CreateAuthTokenRequest

	err := request.DecodeJSON(res, req, &input)
	if err != nil {
		_errors.BadRequest(res, req, err)
		return
	}

	existingUser, err := r.appl.GetByEmailUseCase(input.Email)
	if err != nil && err.Error() != domain.ErrRecordNotFound.Error() {
		_errors.ServerError(res, req, err)
		return
	}

	input.Validator.CheckField(input.Email != "", "Email", "Email is required")
	input.Validator.CheckField(existingUser != nil, "Email", "Email address could not be found")

}
