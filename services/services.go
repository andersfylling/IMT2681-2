package services

import (
	"fmt"
	"time"

	"github.com/andersfylling/IMT2681-2/services/currencyfetcher"
	"github.com/andersfylling/IMT2681-2/services/service"
)

var runnableServices map[string]service.Interface

// Load loads all services into memory and initiate their load func
func Load() {
	runnableServices = make(map[string]service.Interface)

	// Add desired services
	cf := currencyfetcher.NewService()
	cf.Load()
	runnableServices[cf.Title()] = service.Interface(cf)

	// load each service
	// for _, srv := range uninitializedServices {
	// 	runnableServices[srv.Title()] = service.Interface(srv)
	// 	srv.Load()
	// }
}

// Run create a loop that runs every service according to their configuration.
// for some it might be once a day, and some once an hour.
// Must be run as a goroutine!
func Run(done chan<- error, rdy chan<- struct{}, sig <-chan struct{}) {

	Load()

	timeout := time.Now()
	lastRun := time.Now()
	close(rdy)
	for {
		select {
		case <-sig: // stop service loop
			time.Sleep(100 * time.Millisecond)
			fmt.Println("\tServices .. ")

			// empty all services
			for _, srv := range runnableServices {
				srv.Empty()
			}

			done <- nil
			return

		default: // run services
			time.Sleep(time.Millisecond * 500)

			// check if we can run the services
			if !timeout.Before(time.Now().UTC().Add(1 * time.Hour)) {
				continue
			}

			// get a timeout, this way it will never be below a timeout
			// if it is below a timeout it will cause unnecesary func calls/checks
			// A service will always have a lower timeout than 1000 hours.
			nextRun := time.Duration(1000 * time.Hour)

			// For every service initiate their main action
			for _, srv := range runnableServices {
				if srv.Timeout() == time.Duration(0) {
					srv.Run()
				} else if nextRun > srv.Timeout() {
					nextRun = srv.Timeout()
				}
			}
			lastRun = time.Now().UTC().Add(1 * time.Hour)
			timeout = lastRun.Add(nextRun)
		} // select
	} // for
}
