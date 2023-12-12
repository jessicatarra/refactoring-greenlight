package main

import (
	"github.com/jessicatarra/greenlight/internal/config"
	"github.com/jessicatarra/greenlight/internal/errors"
	"net/http"
)

func (a *application) healthcheckHandler(writer http.ResponseWriter, request *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"status":      "available",
			"environment": a.config.Env,
			"version":     config.Version,
		},
	}
	err := a.writeJSON(writer, http.StatusOK, env, nil)
	if err != nil {
		errors.ServerError(writer, request, err)
	}
}
