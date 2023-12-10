package errors

import (
	"github.com/jessicatarra/greenlight/internal/response"
	"github.com/jessicatarra/greenlight/internal/utils/validator"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
)

func ReportError(err error) {
	var (
		message = err.Error()
		trace   = string(debug.Stack())
	)

	slog.Error(message, "trace", trace)
}

func ReportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		method  = r.Method
		url     = r.URL.String()
		trace   = string(debug.Stack())
	)

	requestAttrs := slog.Group("request", "method", method, "url", url)
	slog.Error(message, requestAttrs, "trace", trace)

}

func errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		ReportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerError(w http.ResponseWriter, r *http.Request, err error) {
	ReportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func RateLimitExceeded(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	errorMessage(w, r, http.StatusTooManyRequests, message, nil)
}

//func notFound(w http.ResponseWriter, r *http.Request) {
//	message := "The requested resource could not be found"
//	errorMessage(w, r, http.StatusNotFound, message, nil)
//}

//func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
//	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
//	errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
//}

func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func FailedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		ServerError(w, r, err)
	}
}

func InvalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	errorMessage(w, r, http.StatusUnauthorized, "Invalid authentication token", headers)
}

//func authenticationRequired(w http.ResponseWriter, r *http.Request) {
//	errorMessage(w, r, http.StatusUnauthorized, "You must be authenticated to access this resource", nil)
//}

//
//func basicAuthenticationRequired(w http.ResponseWriter, r *http.Request) {
//	headers := make(http.Header)
//	headers.Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
//
//	message := "You must be authenticated to access this resource"
//	errorMessage(w, r, http.StatusUnauthorized, message, headers)
//}
