package ui

import (
	"net/http"

	"github.com/go-chi/chi"
)

// SetupRoutes configures all the routing done in this application
func setupRoutes(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// fallback
	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("I'm sorry this page does not exist!"))
	})
}
