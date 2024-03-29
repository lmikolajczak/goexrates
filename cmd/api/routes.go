package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	// Initialize a new httprouter router instance.
	router := httprouter.New()

	// Register the relevant methods, basic URL patterns and handler functions for
	// healthcheck and monitoring purposes.
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.Handler(http.MethodGet, "/v1/metrics", expvar.Handler())
	// Application specific routes
	router.HandlerFunc(http.MethodGet, "/v1/latest", app.latestHandler)
	router.HandlerFunc(http.MethodGet, "/v1/historical/:date", app.historicalHandler)

	// Return the httprouter instance (it implements http.Handler interface).
	return app.enableCORS(router)
}
