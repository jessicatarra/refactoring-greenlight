package main

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/database"
	"github.com/jessicatarra/greenlight/internal/validator"
	"net/http"
)

// @Summary Activate User
// @Description Activates a user account using a token that was previously sent when successfully register a new user
// @Tags Users
// @Accept json
// @Produce  json
// @Param token query string true "Token for user activation"
// @Success 200 {object} database.User
// @Router /users/activated [put]
func (app *application) activateUserHandler(writer http.ResponseWriter, request *http.Request) {
	var input struct {
		TokenPlaintext string
	}

	qs := request.URL.Query()

	input.TokenPlaintext = app.readString(qs, "token", "")

	v := validator.New()

	if database.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(writer, request, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(database.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(writer, request, v.Errors)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrEditConflict):
			app.editConflictResponse(writer, request)
		default:
			app.serverErrorResponse(writer, request, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(database.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
		return
	}

	err = app.writeJSON(writer, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(writer, request, err)
	}
}
