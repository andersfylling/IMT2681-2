package currencyfetcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/andersfylling/IMT2681-2/services/service"
)

func getTimeInMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Time since last run
var lastRun = getTimeInMilli()

func getTimeSinceLastRun() float64 {
	return float64(lastRun-getTimeInMilli()) * 1000.0
}

var myClient = &http.Client{Timeout: 5 * time.Second}

func getJSON(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

// Service is of type service.Interface
type Service struct{}

// Info contains the response from http://api.fixer.io/latest?base=EUR
// with json values
type Info struct {
	Base  string         `json:"base"`
	Date  string         `json:"date"`
	Rates map[string]int `json:"rates"`
}

// NewService TODO
func (srv *Service) NewService() service.Interface {
	return &Service{}
}

// Title of the project
func (srv *Service) Title() string {
	return "currencyfetcher"
}

// Load ...
func (srv *Service) Load() {
	// don't need to load anything
}

// Empty ...
func (srv *Service) Empty() {
	// Nothing was loaded into memory
}

// Run ...
func (srv *Service) Run() {
	// create new Info struct
	// get request to the website for getting currency data
	// store the data in the struct or return it?
	fmt.Println("testing service..")

	currency := Info{}
	err := getJSON("https://api.bitfinex.com/v1/pubticker/btcusd", &currency)

	if err != nil {
		fmt.Println(err.Error())
	}
}

// Timeout calculate how much time is left before the next run in seconds
// Will return 0.0 when the service can run
func (srv *Service) Timeout() time.Duration {
	timeout := 2.3 - getTimeSinceLastRun()

	if timeout < 0.0 {
		return 0.0
	}

	return time.Duration(timeout)
}

// make sure struct implements interface
var _ service.Interface = (*Service)(nil)
