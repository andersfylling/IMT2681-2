package notifier

import (
	"fmt"
	"strconv"
	"time"

	"github.com/andersfylling/IMT2681-2/database/dbsession"
	"github.com/andersfylling/IMT2681-2/database/documents/currency"
	"github.com/andersfylling/IMT2681-2/database/documents/webhook"
	"github.com/andersfylling/IMT2681-2/services/currencyfetcher"
	"github.com/andersfylling/IMT2681-2/services/service"
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
	return "notifier"
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
	doc := currency.New()
	doc.Date = today.Format("2006-01-02")
	doc = doc.Find()[0].(*currency.Document) // store results from database

	var results []*webhook.Document

	ses, con, err := dbsession.GetCollection(webhook.Collection)
	if err != nil {
		return
	}
	defer ses.Close()

	// find every web hook
	err = con.Find(nil).All(&results)

	// then check which document should be Invoked
	for _, wh := range results {

		// convert the rate into EUR
		rate := float64(1.0)
		if wh.Base != "EUR" {
			rate = 1 / doc.Rates["EUR"]
		}

		rate *= doc.Rates[wh.Target]

		if wh.Min > rate {
			wh.InvokeWebhook(wh.Target+" dropped below "+strconv.FormatFloat(rate, 'f', -1, 64), "Currency Rate Change!", "")
		} else if wh.Max < rate {
			wh.InvokeWebhook(wh.Target+" went above "+strconv.FormatFloat(rate, 'f', -1, 64), "Currency Rate Change!", "")
		}
	}

	if err != nil {
		fmt.Println(err.Error())
	}
}

// Timeout calculate how much time is left before the next run in seconds
// We want this to run after the currency fetcher
func (srv *Service) Timeout() time.Duration {
	// add 5min to the currency fetchers Timeout
	cf := currencyfetcher.NewService()
	cf.Load()
	return cf.Timeout() + time.Duration(5*time.Minute)
}

// make sure struct implements interface
var _ service.Interface = (*Service)(nil)
