package main

import (
	"net/http"
)

// Declare a handler which writes a response with information about the
// latest exchange rates available in the database.
func (app *application) latestHandler(w http.ResponseWriter, r *http.Request) {
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	codes := app.readCSV(qs, "codes", []string{})

	currencies, date, err := app.models.Currencies.GetLatest(codes)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	rates := map[string]float64{}
	for _, currency := range currencies {
		rates[currency.Code] = currency.Rate
	}
	// Custom JSON response format
	env := envelope{
		"base":  "EUR",
		"date":  date,
		"rates": rates,
	}
	// Send a JSON response containing the currencies data.
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}