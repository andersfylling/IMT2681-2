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

// RateReq is the json struct for incoming request bodies
type RateReq struct {
	Base   string `json:"baseCurrency"`
	Target string `json:"targetCurrency"`
}

// ExchangeRates Responds with the latest exchange rates for given currencies
func ExchangeRates(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	requestedRate := &RateReq{}
	err := decoder.Decode(requestedRate)
	if err != nil {
		fmt.Println(err)
		return
	}

	rate, err := ExchangeRatesFromLatest(w, requestedRate, 0)
	if err != nil {
		w.WriteHeader(503)
		fmt.Println(err.Error())
		return
	}

	w.Write([]byte(strconv.FormatFloat(rate, 'f', -1, 64)))
}

// ExchangeRatesFromLatest Gets the rates for a given date.
// offset can be used to set how many days you want to go from latest or current Date
// It queries the database, and if there exist no entry for given Date
// it queries fixer.io to retrieve it and then stores it into the database.
// the rate is then returned as a float
func ExchangeRatesFromLatest(w http.ResponseWriter, requestedRate *RateReq, offset int) (float64, error) {

	today := time.Now().UTC().Add(1*time.Hour).AddDate(0, 0, -offset) // CET = UTC+1
	doc := currency.New()
	doc.Date = today.Format("2006-01-02")

	// make sure it exists in database, otherwise use fixer.io
	if len(doc.Find()) == 0 {
		// get data from fixer, and save it to database
		err := utils.GetJSON("http://api.fixer.io/"+doc.Date+"?base=EUR", &doc)
		if err != nil {
			return -1.0, err
		}

		// Store data to database
		doc.Insert()
	} else {
		doc = doc.Find()[0].(*currency.Document) // populate with the db data
	}

	// convert the rate into EUR
	rate := float64(1.0)
	if requestedRate.Base != "EUR" {
		rate = 1 / doc.Rates["EUR"]
	}

	rate *= doc.Rates[requestedRate.Target]

	return rate, nil
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
