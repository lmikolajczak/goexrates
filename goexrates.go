package main

import (
	"net/http"
	"os"

	"goexrates/controllers"
	"goexrates/models"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	models.InitDB("postgres://luq:jogabonito13@localhost/goexrates?sslmode=disable")

	router := mux.NewRouter()

	router.HandleFunc("/latest", controllers.Latest)
	router.HandleFunc("/{date}", controllers.Historical)

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)

	http.ListenAndServe(":3000", loggedRouter)
}
