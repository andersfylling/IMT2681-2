package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/andersfylling/IMT2681-2/database"
	"github.com/andersfylling/IMT2681-2/services"
	"github.com/andersfylling/IMT2681-2/ui"
)

func main() {
	termSignal := make(chan os.Signal, 1)
	stop := make(chan struct{})
	signal.Notify(termSignal, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	fmt.Println("Starting up")

	// start up database
	databaseChan := make(chan error)
	databaseRdyChan := make(chan struct{})
	fmt.Print("\tConnecting to database .. ")
	go database.Connect(databaseChan, databaseRdyChan, stop)

	// start web server
	<-databaseRdyChan
	fmt.Println("OK")
	webserverChan := make(chan error)
	webserverRdyChan := make(chan struct{})
	fmt.Print("\tStarting web server .. ")
	go ui.UI(webserverChan, webserverRdyChan, stop)

	// start services
	<-webserverRdyChan
	fmt.Println("OK")
	servicesChan := make(chan error)
	servicesRdyChan := make(chan struct{})
	fmt.Print("\tInitiating services .. ")
	go services.Run(servicesChan, servicesRdyChan, stop)

	<-servicesRdyChan
	fmt.Println("OK")
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
