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
