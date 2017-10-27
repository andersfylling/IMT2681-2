package average

import "net/http"

// ForLastSevenDays Responds with the avg for the last seven days
func ForLastSevenDays(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Average for last seven days"))
}

// Info some details about this uri
func Info(w http.ResponseWriter, r *http.Request) {
	res := ""
	res += "<pre>\n"
	res += "[POST] /average\n"
	res += "\tGet the average exchange rates given a base and target currency\n"
	res += "\tCurrencies are given by a json body:\n"
	res += "<code>\n"
	res += "{\n\t\"baseCurrency\":\"EUR\",\n\t\"targetCurrency\":\"NOK\"\n}\n"
	res += "</code></pre>"

	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(res))
}
