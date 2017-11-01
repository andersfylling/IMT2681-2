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
	// Add desired services
	uninitializedServices := []service.Interface{
		(*currencyfetcher.Service)(nil),
	}

	// load each service
	for _, srv := range uninitializedServices {
		runnableServices[srv.Title()] = srv
		srv.Load()
	}
}

// Run create a loop that runs every service according to their configuration.
// for some it might be once a day, and some once an hour.
// Must be run as a goroutine!
func Run(done chan<- error, rdy chan<- struct{}, sig <-chan struct{}) {

	close(rdy)
	for {
		timeout := time.Duration(0)
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
			time.Sleep(timeout * time.Millisecond)
			// get a timeout, this way it will never be below a timeout
			// if it is below a timeout it will cause unnecesary func calls/checks
			for _, srv := range runnableServices {
				timeout = srv.Timeout()
				break // This is the ugliest way I've seen in my life to access a
				// random/first entry
			}

			// For every service initiate their main action
			for _, srv := range runnableServices {
				if srv.Timeout() == 0 {
					srv.Run()
				} else if timeout > srv.Timeout() {
					timeout = srv.Timeout()
				}
			}
		}
	}
}
