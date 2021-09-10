package main

import (
	"net/http"
	"time"
)

// Declare a handler which writes a response with information about the
// latest exchange rates available in the database.
func (app *application) latestHandler(w http.ResponseWriter, r *http.Request) {
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	codes := app.readCSV(qs, "codes", []string{})
	base := app.readString(qs, "base", "EUR")

	rates, err := app.models.Currencies.GetRates(time.Time{}, codes)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if base != "EUR" {
		app.models.Currencies.Convert(base, rates)
	}

	serializer := RatesSerializer{data: rates, base: base}
	err = app.writeJSON(w, http.StatusOK, serializer.dump(), nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
