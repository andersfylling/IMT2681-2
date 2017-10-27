package webhook

import "net/http"

// CreateWebhook ..
func CreateWebhook(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Register webhook"))
}

// Info some details about this uri
func Info(w http.ResponseWriter, r *http.Request) {
	res := ""
	res += "<pre>\n"
	res += "[POST] /:\n"
	res += "\tNew webhooks can be registered using POST requests with the\n"
	res += "\tfollowing schema. Note we will use /root as a placeholder for\n"
	res += "\tthe root path of your web service (i.e. the path you will submit\n"
	res += "\tto the submission system). For example, if your web service runs\n"
	res += "\ton https://localhost:8080/exchange, then this is the root path\n"
	res += "\tyou would submit.\n"
	res += "<code>\n"
	res += "{\n"
	res += "\t\"webhookURL\": \"http://remoteUrl:8080/randomWebhookPath\",\n"
	res += "\t\"baseCurrency\": \"EUR\",\n"
	res += "\t\"targetCurrency\": \"NOK\",\n"
	res += "\t\"minTriggerValue\": 1.50,\n"
	res += "\t\"maxTriggerValue\": 2.55\n"
	res += "}\n"
	res += "</code></pre>"

	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(res))
}
