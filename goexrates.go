package main

import (
	"net/http"

	"github.com/Luqqk/goexrates/handlers"
	"github.com/Luqqk/goexrates/models"
	"github.com/gorilla/mux"
)

func main() {
	models.InitDB("postgres://luq:jogabonito13@localhost/goexrates?sslmode=disable")
	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/latest", handlers.Latest)
	router.HandleFunc("/{date}", handlers.Historical)

	http.ListenAndServe(":3000", router)
}
