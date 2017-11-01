package routes

import (
	"net/http"
	"net/url"

	"github.com/andersfylling/IMT2681-2/ui/routes/average"
	"github.com/andersfylling/IMT2681-2/ui/routes/evaluationtrigger"
	"github.com/andersfylling/IMT2681-2/ui/routes/latest"
	"github.com/andersfylling/IMT2681-2/ui/routes/webhook"
	"github.com/go-chi/chi"
)

// https://golangcode.com/how-to-check-if-a-string-is-a-url/
func validURL(link string) bool {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return false
	}

	return true
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("I'm sorry this page does not exist!"))
}

// SetupRoutes configures all the routing done in this application
func SetupRoutes(router *chi.Mux) {
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	// webhooks
	router.Post("/", webhook.CreateWebhook)
	router.Post("/webhookUrl", webhook.InvokeWebhook)
	router.Get("/", webhook.Info)

	// latest exchange rates
	router.Post("/latest", latest.ExchangeRates)
	router.Get("/latest", latest.Info)

	// average exchange rates
	router.Post("/average", average.ForLastSevenDays)
	router.Get("/average", average.Info)

	// average exchange rates
	router.Post("/evaluationtrigger", evaluationtrigger.FireAllWebhooks)
	router.Get("/evaluationtrigger", evaluationtrigger.Info)

	// invoking a webhook, otherwise a fallback
	router.Get("/{webhookURL}", func(w http.ResponseWriter, r *http.Request) {
		webhookURL := chi.URLParam(r, "webhookURL")

		if validURL(webhookURL) {
			// check if it exist in database
			// if it does not exist in the database call notFound(...)
		} else {
			notFound(w, r)
		}
	})
}
