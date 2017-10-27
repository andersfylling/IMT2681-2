package latest

import "net/http"

// ExchangeRates Responds with the latest exchange rates for given currencies
func ExchangeRates(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("here you get an latest exchange rates"))
}

// Info some details about this uri
func Info(w http.ResponseWriter, r *http.Request) {
	res := ""
	res += "<pre>\n"
	res += "[POST] /latest:\n"
	res += "\tGet the latest exchange rates given a base and target currency\n"
	res += "\tCurrencies are given by a json body:\n"
	res += "<code>\n"
	res += "{\n\t\"baseCurrency\":\"EUR\",\n\t\"targetCurrency\":\"NOK\"\n}\n"
	res += "</code></pre>"

	w.Header().Set("Content-Type", "text/html")

	w.Write([]byte(res))
}
