package main

import (
	"net/http"

	"github.com/Luqqk/goexrates/handlers"
	"github.com/Luqqk/goexrates/middlewares"
	"github.com/Luqqk/goexrates/models"
	"github.com/gorilla/mux"
)

func main() {
	models.InitDB("postgres://luq:jogabonito13@localhost/goexrates?sslmode=disable")
	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/latest", middlewares.CORS(handlers.Latest))
	router.HandleFunc("/{date}", middlewares.CORS(handlers.Historical))

	http.ListenAndServe(":3000", middlewares.Log(router))
}
