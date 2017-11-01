package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	webhookDoc "github.com/andersfylling/IMT2681-2/database/documents/webhook"
	"github.com/andersfylling/IMT2681-2/ui/routes/average"
	"github.com/andersfylling/IMT2681-2/ui/routes/evaluationtrigger"
	"github.com/andersfylling/IMT2681-2/ui/routes/latest"
	"github.com/andersfylling/IMT2681-2/ui/routes/webhook"
	"github.com/go-chi/chi"
	"gopkg.in/mgo.v2/bson"
)

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
	router.Get("/{webhookID}", func(w http.ResponseWriter, r *http.Request) {
		webhookID := chi.URLParam(r, "webhookID")

		wh := webhookDoc.New()
		wh.ID = bson.ObjectIdHex(webhookID)

		if len(wh.Find()) > 0 {
			jsonStr, err := json.Marshal(wh.Find()[0].(*webhookDoc.Document))
			if err != nil {
				w.WriteHeader(503)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(jsonStr))
			w.WriteHeader(200)
			return
		} else {
			notFound(w, r)
		}
	})
	router.Delete("/{webhookID}", func(w http.ResponseWriter, r *http.Request) {
		webhookID := chi.URLParam(r, "webhookID")

		wh := webhookDoc.New()
		wh.ID = bson.ObjectIdHex(webhookID)

		if len(wh.Find()) > 0 {
			deleted := len(wh.Remove())
			if deleted == 0 {
				w.WriteHeader(503)
			} else {
				w.WriteHeader(200)
			}
		} else {
			w.WriteHeader(204)
		}
	})
}
