package helpers

import (
	"errors"
	"github.com/jessicatarra/greenlight/internal/utils/validator"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"net/url"
	"strconv"
)

type Helpers interface {
	ReadIDParam(request *http.Request) (int64, error)
	ReadString(qs url.Values, key string, defaultValue string) string
	ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int
}

type helpers struct{}

func New() Helpers {
	return &helpers{}
}

func (h *helpers) ReadIDParam(request *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(request.Context())

	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

func (h *helpers) ReadString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (h *helpers) ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError("must be an integer value")
		return defaultValue
	}

	return i
}
