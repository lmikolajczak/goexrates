package middlewares

import (
	"log"
	"net/http"
)

// Log server requests
func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("--> %s %s", r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
