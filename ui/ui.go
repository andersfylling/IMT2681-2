package ui

import (
	"context"
	"fmt"
	"os"
	"time"

	"net/http"

	"github.com/andersfylling/IMT2681-2/ui/middlewares"
	"github.com/andersfylling/IMT2681-2/ui/routes"
	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
)

// UI Set up a web interface in a seperate thread
// Inspired from: https://github.com/btcsuite/btcd/blob/master/btcd.go
func UI(done chan error, sig chan struct{}) {
	fmt.Println("Starting http server..")

	// Our graceful valve shut-off package to manage code preemption and
	// shutdown signaling.
	valv := valve.New()
	baseCtx := valv.Context()

	// HTTP service running in this program as well. The valve context is set
	// as a base context on the server listener at the point where we instantiate
	// the server - look lower.
	router := chi.NewRouter()
	middlewares.Setup(router)
	routes.SetupRoutes(router)

	// setup port to listen to
	var port string
	if os.Getenv("PORT") == "" {
		port = ":8080" // if the heroku port is not specified, run on port 8080
	} else {
		port = ":" + os.Getenv("PORT")
	}

	srv := http.Server{Addr: port, Handler: chi.ServerBaseContext(baseCtx, router)}

	go func() {
		for {
			select {
			case <-sig:
				// time to terminate web server
				fmt.Println("\tWeb server .. ")

				// first valv
				valv.Shutdown(1 * time.Second)

				// create context with timeout
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				// start http shutdown
				srv.Shutdown(ctx)

				// verify, in worst case call cancel via defer
				select {
				case <-time.After(5 * time.Second):
					fmt.Println("not all connections done")
				case <-ctx.Done():
					return
				}
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()
	srv.ListenAndServe()

	done <- nil
}
