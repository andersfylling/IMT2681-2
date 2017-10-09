package services

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/valve"
	"github.com/sciencefyll/IMT2681-2/services/currencyfetcher"
)

var runnableServices map[string]ServiceInterface

func Load() {
	runnableServices["currencyfetcher"] = currencyfetcher.Info
}

// Run create a loop that runs every service according to their configuration.
// for some it might be once a day, and some once an hour.
func Run(done chan error) {
	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	valv := valve.New()
	//baseCtx := valv.Context()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			fmt.Println("shutting down services..")

			// first valv
			valv.Shutdown(20 * time.Second)

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()

			// start http shutdown
			// ....

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(21 * time.Second):
				fmt.Println("not all connections done")
			case <-ctx.Done():

			}
		}
	}()

	done <- nil
}
