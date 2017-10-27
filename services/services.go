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
func Run(done chan error, sig chan struct{}) {
	fmt.Println("Starting service loop") // use logger for printing..

	// the run loop
	go func() {
		for {
			select {
			case <-sig: // stop service loop
				time.Sleep(100 * time.Millisecond)
				fmt.Print("\tServices .. ")

				// empty all services
				for _, srv := range runnableServices {
					srv.Empty()
				}

				fmt.Println("OK")
				done <- nil

				return

			default: // run services
				for _, srv := range runnableServices {
					srv.Run()
				}
			}
		}
	}()
}
