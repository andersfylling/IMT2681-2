package ui

import (
	"context"
	"fmt"
	"os"
	"time"

	"net/http"

	"github.com/andersfylling/IMT2681-2/ui/middlewares"
	"github.com/go-chi/chi"
	"github.com/go-chi/valve"
)

var baseURI = "/"

// Configure sets the base uri for the web server
// @param baseUri string does not need to start nor end with `/`
func Configure(baseUri string) {
	baseURI = baseUri

	// make sure there are forward slashes
	baseURI = "/" + baseURI + "/"

	// trim duplicate slashes
	var result string
	previous := '¤' // assumption: a uri never starts with ¤
	for _, r := range baseURI {
		if r == '/' && r != previous {
			result += string(r)
		}
		previous = r
	}

	// update the base uri
	baseURI = result
}

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

	router.Get(baseURI, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sup"))
	})

	router.Get(baseURI+"slow", func(w http.ResponseWriter, r *http.Request) {

		valve.Lever(r.Context()).Open()
		defer valve.Lever(r.Context()).Close()

		select {
		case <-valve.Lever(r.Context()).Stop():
			fmt.Println("valve is closed. finish up..")

		case <-time.After(5 * time.Second):
			// The above channel simulates some hard work.
			// We want this handler to complete successfully during a shutdown signal,
			// so consider the work here as some background routine to fetch a long running
			// search query to find as many results as possible, but, instead we cut it short
			// and respond with what we have so far. How a shutdown is handled is entirely
			// up to the developer, as some code blocks are preemptable, and others are not.
			time.Sleep(5 * time.Second)
		}

		w.Write([]byte(fmt.Sprintf("all done.\n")))
	})

	// setup port to listen to
	var port string
	if os.Getenv("PORT") == "" {
		port = ":8080" // if the heroku port is not specified, run on port 8080
	} else {
		port = ":" + os.Getenv("PORT")
	}

	srv := http.Server{Addr: port, Handler: chi.ServerBaseContext(baseCtx, router)}

	go func() {
		for range sig {
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
				fmt.Println("OK")
				return
			}
		}
	}()
	srv.ListenAndServe()

	done <- nil
}
