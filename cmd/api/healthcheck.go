package main

import (
	"net/http"
)

// Declare a handler which writes a response with information about the
// application status, operating environment and version.
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	serializer := HealthCheckSerializer{data: map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}}
	err := app.writeJSON(w, http.StatusOK, serializer.dump(), nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
