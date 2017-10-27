package evaluationtrigger

import "net/http"

// FireAllWebhooks Fires all webhooks
func FireAllWebhooks(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("here you get an latest exchange rates"))
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
