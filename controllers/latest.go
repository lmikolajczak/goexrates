package controllers

import (
	"encoding/json"
	"goexrates/models"
	"net/http"
)

// Latest /latest route controller
func Latest(w http.ResponseWriter, r *http.Request) {
	baseParam := r.URL.Query().Get("base")
	symbolsParam := r.URL.Query().Get("symbols")

	currencies, err := models.LatestRates(baseParam, symbolsParam)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	response, err := json.Marshal(currencies)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
