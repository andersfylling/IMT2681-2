package latest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andersfylling/IMT2681-2/database/documents/currency"
	"github.com/andersfylling/IMT2681-2/utils"
)

type RateReq struct {
	Base   string `json:"baseCurrency"`
	Target string `json:"targetCurrency"`
}

// ExchangeRates Responds with the latest exchange rates for given currencies
func ExchangeRates(w http.ResponseWriter, r *http.Request) {

	today := time.Now().UTC().Add(1 * time.Hour) // CET = UTC+1
	doc := currency.New()
	doc.Date = today.Format("2006-01-02")

	// make sure it exists in database, otherwise use fixer.io
	if len(doc.Find()) == 0 {
		// get data from fixer, and save it to database
		err := utils.GetJSON("http://api.fixer.io/latest?base=EUR", &doc)
		if err != nil {
			fmt.Println(err.Error())
			w.WriteHeader(503)
			return
		}

		// Store data to database
		doc.Insert()
	}

	decoder := json.NewDecoder(r.Body)
	requestedRate := &RateReq{}
	err := decoder.Decode(requestedRate)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(503)
		return
	}

	// convert the rate into EUR
	rate := float64(1.0)
	if requestedRate.Base != "EUR" {
		rate = 1 / doc.Rates["EUR"]
	}

	rate *= doc.Rates[requestedRate.Target]

	w.Write([]byte(strconv.FormatFloat(rate, 'f', -1, 64)))
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
