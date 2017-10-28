package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/andersfylling/IMT2681-2/database"
	"github.com/andersfylling/IMT2681-2/services"
	"github.com/andersfylling/IMT2681-2/ui"
)

func main() {
	termSignal := make(chan os.Signal, 1)
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	stop := make(chan struct{})

	// start up database
	databaseChan := make(chan error)
	go database.Connect(databaseChan, stop)

	// start web server
	webserverChan := make(chan error)
	go ui.UI(webserverChan, stop)

	// start services
	servicesChan := make(chan error)
	go services.Run(servicesChan, stop)

	time.Sleep(10 * time.Millisecond)
	fmt.Println("Program is now running.  Press CTRL-C to exit.")

	// Shut down using timeout by simulating the interupt signal
	// go func() {
	// 	time.Sleep(3 * time.Second)
	//
	// 	termSignal <- os.Signal(os.Interrupt)
	// }()

	<-termSignal
	fmt.Println("\nShutting down..")
	close(stop)

	// wait for things to completely stop
	<-servicesChan
	fmt.Println("\tServices OK")
	<-webserverChan
	fmt.Println("\tWeb server OK")
	<-databaseChan
	fmt.Println("\tDatabase OK")

	// all is done
	fmt.Println("Shutdown complete")
}
