package currencyfetcher

import (
	"fmt"
	"time"

	"github.com/andersfylling/IMT2681-2/database/documents/currency"
	"github.com/andersfylling/IMT2681-2/services/service"
	"github.com/andersfylling/IMT2681-2/utils"
)

var lastRun time.Duration
var firstRun bool

func getTimeInMilli() time.Duration {
	return time.Duration(time.Now().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond)))
}
func getTimeSinceLastRun() time.Duration {
	return getTimeInMilli() - lastRun
}

// Time since last run
func init() {
	lastRun = 0 // So that it runs as soon as the service starts up
	firstRun = true
}

// Service is of type service.Interface
type Service struct{}

// NewService TODO
func NewService() *Service {
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
	lastRun = getTimeInMilli()

	t := time.Now()
	today := t.UTC().Add(1 * time.Hour) // CET = UTC+1
	currency := currency.New()
	currency.Date = today.Format("2006-01-02")

	// make sure we haven't already a rate for today
	if len(currency.Find()) > 0 {
		return // already have a currency rate
	}

	err := utils.GetJSON("http://api.fixer.io/latest?base=EUR", &currency)
	if err != nil {
		fmt.Println(err.Error())
	}

	// Store data to database
	currency.Insert()
}

// Timeout calculate how much time is left before the next run in seconds
// Will return 0.0 when the service can run
func (srv *Service) Timeout() time.Duration {
	last := getTimeSinceLastRun()

	// Get the time until the 4pm cet time
	// if the `until` has passed, say this is run at 5pm
	// return 0, so it can be run at once.

	// however, if the app gets started before 4pm
	// we want this to run once to get the rates.
	if firstRun == true {
		firstRun = false
		return time.Duration(0)
	}

	// in case its above 24hours since last run, just execute
	if last/time.Hour > 24 {
		return time.Duration(0)
	}

	t := time.Now().UTC().Add(1 * time.Hour) // CET = UTC+1

	if t.Hour() < 16 {
		return time.Duration(16-t.Hour())*time.Hour - (time.Duration(t.Minute()) * time.Minute)
	} else if t.Hour() == 16 && last/time.Hour > 1 {
		return time.Duration(0)
	}

	return time.Duration(16+(24-t.Hour()))*time.Hour - (time.Duration(t.Minute()) * time.Minute)
}

// make sure struct implements interface
var _ service.Interface = (*Service)(nil)
