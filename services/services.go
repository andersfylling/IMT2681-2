package services

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/andersfylling/IMT2681-2/services/currencyfetcher"
	"github.com/andersfylling/IMT2681-2/services/service"
	"github.com/go-chi/valve"
)

var runningSrvsChan chan struct{}
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

// Stop shuts down services
func Stop() {
	// stop the services
	fmt.Println("shutting down services..")

	// make sure the service loop is stopped
	close(runningSrvsChan)

	// empty all services
	for _, srv := range runnableServices {
		srv.Empty()
	}
}

// Run create a loop that runs every service according to their configuration.
// for some it might be once a day, and some once an hour.
func Run(done chan error) {
	// let the prog know it can run
	runningSrvsChan = make(chan struct{}) // to end the forever service loop

	fmt.Println("Starting service loop") // use logger for printing..

	// the run loop
	go func() {
		for {
			select {
			case <-runningSrvsChan:
				fmt.Println("\tStopped the service loop")
				return
			default:
				for _, srv := range runnableServices {
					srv.Run()
				}
			}
		}
	}()

	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	valv := valve.New()
	//baseCtx := valv.Context()

	// TODO: is this really needed? some, but not the last part of the loop
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			Stop()

			// first valv
			valv.Shutdown(3 * time.Second) // 3s timeout

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			// start http shutdown
			// ....

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(5 * time.Second):
				fmt.Println("not all services were shut down")
			case <-ctx.Done():

			}
		}
	}()

	done <- nil
}
