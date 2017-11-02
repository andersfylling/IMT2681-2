package average

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestForLastSevenDays(t *testing.T) {
	var data = []byte(`{"baseCurrency": "EUR","targetCurrency":"NOK"}`)
	r, err := http.NewRequest("POST", "/average", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ForLastSevenDays)
	handler.ServeHTTP(rr, r)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Check the response body is what we expect.
	f, err := strconv.ParseFloat(rr.Body.String(), 64)
	if err != nil {
		t.Errorf("Unable to convert body to float: got %v", rr.Body.String())
	}
	if f < 9.0 {
		t.Errorf("Making sure the average is above 9nok/eur: got %v want %v", f, 9.0)
	}
}
