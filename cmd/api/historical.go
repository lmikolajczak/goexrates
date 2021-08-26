package main

import (
	"net/http"
)

// Declare a handler which writes a response with information about the
// historical exchange rates available in the database for specific date.
func (app *application) historicalHandler(w http.ResponseWriter, r *http.Request) {
	date, err := app.readDateParam(r)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	// Custom JSON response format
	env := envelope{
		"status": "not implemented",
		"date":   date,
	}
	// Send a JSON response containing the currencies data.
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
