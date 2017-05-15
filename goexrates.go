package main

import (
	"encoding/json"
	"goexrates/models"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	models.InitDB("postgres://user:password@localhost/dbname?sslmode=disable")

	router := mux.NewRouter()

	router.HandleFunc("/latest", Latest)
	router.HandleFunc("/{date}", Historical)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	http.ListenAndServe(":3000", loggedRouter)
}

// Latest : /latest route handler
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

// Historical : /{date} route handler
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
