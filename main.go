package main

import (
	"fmt"

	"github.com/andersfylling/IMT2681-2/services"
	"github.com/andersfylling/IMT2681-2/ui"
)

func main() {
	// start web server
	webserverChan := make(chan error)
	go ui.UI(webserverChan)

	// start services
	servicesChan := make(chan error)
	go services.Run(servicesChan)

	// wait for shutdown signals
	<-webserverChan
	<-servicesChan

	// all is done
	fmt.Println("DONE")
}
