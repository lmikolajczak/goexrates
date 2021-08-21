package main

import (
	"net/http"
)

// Declare a handler which writes a response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		http.Error(
			w, "The server encountered a problem and could not process your request",
			http.StatusInternalServerError,
		)
	}
}
