package main

import (
	"encoding/json"
	"goexrates/models"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	models.InitDB("postgres://user:password@localhost/dbname?sslmode=disable")
	router := mux.NewRouter()
	router.HandleFunc("/latest", Latest)
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	http.ListenAndServe(":3000", loggedRouter)
}

//Latest : /latest route handler
func Latest(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	symbols := r.URL.Query().Get("symbols")
	if base == "EUR" {
		base = ""
	}

	var currencies *models.Currencies
	var err error

	switch {
	case base != "" && symbols != "":
		currencies, err = models.RecalculatedAndFilteredRates(base, symbols)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	case base != "" && symbols == "":
		currencies, err = models.RecalculatedRates(base)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	case base == "" && symbols != "":
		currencies, err = models.FilteredRates(symbols)
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	default:
		currencies, err = models.LatestRates()
		if err != nil {
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}

	resp, err := json.Marshal(currencies)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
