package http

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/database"
	_errors "github.com/jessicatarra/greenlight/internal/errors"
	"github.com/jessicatarra/greenlight/internal/password"
	"github.com/jessicatarra/greenlight/internal/request"
	"github.com/jessicatarra/greenlight/internal/response"
	"github.com/jessicatarra/greenlight/internal/utils/helpers"
	"github.com/jessicatarra/greenlight/ms/auth/internal/domain"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type envelope map[string]interface{}

type Handlers interface {
	createUser(res http.ResponseWriter, req *http.Request)
	activateUser(res http.ResponseWriter, req *http.Request)
	createAuthenticationToken(res http.ResponseWriter, req *http.Request)
}

type handlers struct {
	appl    domain.Appl
	helpers helpers.Helpers
}

func (s service) Handlers(router *httprouter.Router) {
	res := registerHandlers(s.appl)

	router.HandlerFunc(http.MethodPost, "/v1/users", res.createUser)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", res.activateUser)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", res.createAuthenticationToken)
}

func registerHandlers(appl domain.Appl) Handlers {
	return &handlers{
		appl:    appl,
		helpers: helpers.New(),
	}
}

// @Summary Register User
// @Description Registers a new user.
// @Tags Users
// @Accept json
// @Produce  json
// @Param name body domain.CreateUserRequest true "User registration data"
// @Success 201 {object} domain.User
// @Router /users [post]
func (h *handlers) createUser(res http.ResponseWriter, req *http.Request) {
	var input domain.CreateUserRequest

	err := request.DecodeJSON(res, req, &input)
	if err != nil {
		_errors.BadRequest(res, req, err)
		return
	}

	existingUser, err := h.appl.GetByEmailUseCase(input.Email)
	if err != nil && err.Error() != domain.ErrRecordNotFound.Error() {
		_errors.ServerError(res, req, err)
		return
	}

	ValidateUser(input, existingUser)

	if input.Validator.HasErrors() {
		_errors.FailedValidation(res, req, input.Validator)
		return
	}

	hashedPassword, err := password.Hash(input.Password)
	if err != nil {
		_errors.ServerError(res, req, err)
		return
	}

	user, err := h.appl.CreateUseCase(input, hashedPassword)
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
func (h *handlers) activateUser(res http.ResponseWriter, req *http.Request) {
	var input domain.ActivateUserRequest

	qs := req.URL.Query()

	input.TokenPlaintext = h.helpers.ReadString(qs, "token", "")

	ValidateToken(input)

	if input.Validator.HasErrors() {
		_errors.FailedValidation(res, req, input.Validator)
		return
	}

	user, err := h.appl.ActivateUseCase(input.TokenPlaintext)
	if err != nil {
		_errors.ServerError(res, req, err)
		return
	}

	err = response.JSON(res, http.StatusCreated, envelope{"user": user})
	if err != nil {
		_errors.ServerError(res, req, err)
	}
}

// @Summary Create authentication token
// @Description Creates an authentication token for a user
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body domain.CreateAuthTokenRequest true "Request body"
// @Success 201 {object} domain.Token "Authentication token"
// @Router /tokens/authentication [post]
func (h *handlers) createAuthenticationToken(res http.ResponseWriter, req *http.Request) {
	var input domain.CreateAuthTokenRequest

	err := request.DecodeJSON(res, req, &input)
	if err != nil {
		_errors.BadRequest(res, req, err)
		return
	}

	existingUser, err := h.appl.GetByEmailUseCase(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrRecordNotFound):
			_errors.InvalidAuthenticationToken(res, req)
		default:
			_errors.ServerError(res, req, err)
		}
		return
	}

	ValidateEmailForAuth(input, existingUser)

	if existingUser != nil {
		passwordMatches, err := password.Matches(input.Password, existingUser.HashedPassword)
		if err != nil {
			_errors.ServerError(res, req, err)
			return
		}

		ValidatePasswordForAuth(input, passwordMatches)
	}

	if input.Validator.HasErrors() {
		_errors.FailedValidation(res, req, input.Validator)
		return
	}

	jwtBytes, err := h.appl.CreateAuthTokenUseCase(existingUser.ID)
	if err != nil {
		_errors.ServerError(res, req, err)
		return
	}

	err = response.JSON(res, http.StatusCreated, envelope{"authentication_token": string(jwtBytes)})
	if err != nil {
		_errors.ServerError(res, req, err)
	}

}
