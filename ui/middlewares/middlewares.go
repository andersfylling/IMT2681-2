package middlewares

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func removeDuplicateSlashInURI(uri string) string {
	product := ""
	var previous rune = '¤'
	for _, r := range uri {
		if !(previous == r && r == '/') {
			product += string(r)
		}

		previous = r
	}
	return product
}

// Setup smth
func Setup(router *chi.Mux) {

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)

	// remove duplicates of slashes in URI: /// => /
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.RequestURI = removeDuplicateSlashInURI(r.RequestURI)
			next.ServeHTTP(w, r)
		})
	})
}
