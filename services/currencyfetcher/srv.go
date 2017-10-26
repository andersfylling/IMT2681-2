package currencyfetcher

import (
	"fmt"
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
}

// Timeout calculate how much time is left before the next run in seconds
// Will return 0.0 when the service can run
func (srv *Service) Timeout() float64 {
	timeout := 2.3 - getTimeSinceLastRun()

	if timeout < 0.0 {
		return 0.0
	}

	return timeout
}

// make sure struct implements interface
var _ service.Interface = (*Service)(nil)
