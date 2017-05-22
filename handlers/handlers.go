package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Luqqk/goexrates/models"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
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

	decimal.MarshalJSONWithoutQuotes = true
	response, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

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

	decimal.MarshalJSONWithoutQuotes = true
	response, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

// Favicon /favicon.ico route controller
func Favicon(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	http.ServeFile(w, r, "./static/favicon.ico")
}
