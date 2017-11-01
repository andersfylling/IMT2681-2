package evaluationtrigger

import (
	"net/http"

	"github.com/andersfylling/IMT2681-2/database/documents/webhook"
)

// FireAllWebhooks Fires all webhooks
func FireAllWebhooks(w http.ResponseWriter, r *http.Request) {
	arr, err := webhook.InvokeAll()
	if err != nil {
		w.WriteHeader(503)
		return
	}

	// if there was a match send OK
	if len(arr) > 0 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(204) // No content
	}

	return
}

// Info some details about this uri
func Info(w http.ResponseWriter, r *http.Request) {
	res := ""
	res += "<pre>\n"
	res += "[GET] /evaluationtrigger:\n"
	res += "\tThis request invokes all webhooks (i.e. bypasses the timed trigger)\n"
	res += "\tand sends the payload as specified under 'Invoking a registered"
	res += "webhook'.\n"
	res += "\tThis functionality is meant for testing and evaluation purposes.\n"
	res += "</pre>"

	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(res))
}
