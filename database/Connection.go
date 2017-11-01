package database

import (
	"fmt"
	"time"

	"github.com/andersfylling/IMT2681-2/database/dbsession"
)

// Connect Must be run as a goroutine
func Connect(done chan<- error, rdy chan<- struct{}, sig <-chan struct{}) {
	session, err := dbsession.GetInstance()
	if err == nil {
		err = session.Ping()
	}
	if err != nil {
		fmt.Println(fmt.Errorf("DATABASE-ERROR: " + err.Error()))
		close(rdy)
		done <- err
		return
	}
	close(rdy)

	<-sig // stop service loop
	time.Sleep(100 * time.Millisecond)
	fmt.Println("\tDatabase session .. ")

	err = session.Ping() // maybe use refresh then close?

	if err != nil {
		done <- nil // error cant pint closed connection
		return
	}

	// Ping worked, so close the connection
	session.Close()
	done <- nil
}
