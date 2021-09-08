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
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	codes := app.readCSV(qs, "codes", []string{})

	currencies, resultsDate, err := app.models.Currencies.GetRates(date, codes)
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
		"source": "EUR",
		"date":   resultsDate,
		"rates":  rates,
	}
	// Send a JSON response containing the currencies data.
	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
