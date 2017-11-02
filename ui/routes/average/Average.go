package average

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/andersfylling/IMT2681-2/ui/routes/latest"
)

// ForLastSevenDays Responds with the avg for the last seven days
func ForLastSevenDays(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	requestedRate := &latest.RateReq{}
	err := decoder.Decode(requestedRate)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(503)
		return
	}

	// create a check for every day and store every rate thats missing from the db
	sum := float64(0.0)
	retries := 0
	for offset := 0; offset < 7 && retries < 3; offset++ {
		rate, err := latest.ExchangeRatesFromLatest(w, requestedRate, offset)
		if err != nil {
			fmt.Println("loop issue", err.Error())
			w.WriteHeader(503)
			return
		}

		if rate == 0.0 {
			offset--
			retries++
			time.Sleep(150 * time.Millisecond)
		}

		sum += rate
	}

	if retries == 3 {
		w.WriteHeader(408)
	} else {
		w.WriteHeader(200)
		w.Write([]byte(strconv.FormatFloat(sum/7, 'f', -1, 64)))
	}
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
