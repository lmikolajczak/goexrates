package controllers

import (
	"encoding/json"
	"goexrates/models"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// Historical : /{date} route controller
func Historical(w http.ResponseWriter, r *http.Request) {
	baseParam := r.URL.Query().Get("base")
	symbolsParam := r.URL.Query().Get("symbols")
	vars := mux.Vars(r)

	_, err := time.Parse("2006-01-02", vars["date"])
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	currencies, err := models.HistoricalRates(baseParam, symbolsParam, vars["date"])
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
