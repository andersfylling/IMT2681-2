package database

import (
	"errors"
	"fmt"
	"os"

	"sync"

	mgo "gopkg.in/mgo.v2"
)

type singleton struct {
	session *mgo.Session
}

var instance *singleton
var once sync.Once

// GetInstance singleton pattern
func GetInstance() (*mgo.Session, error) {
	var issue error = nil
	once.Do(func() {
		instance = &singleton{}

		// get username
		if os.Getenv("MGO_USER") == "" {
			issue = errors.New("mgo: missing username in environment var MGO_USER")
			return
		}
		// get password
		if os.Getenv("MGO_PASS") == "" {
			issue = errors.New("mgo: missing password in environment var MGO_PASS")
			return
		}

		// store them
		user := os.Getenv("MGO_USERNAME")
		pass := os.Getenv("MGO_PASS")

		// connect to external database
		ses, err := mgo.Dial("mongodb://<" + user + ">:<" + pass + ">@ds045099.mlab.com:45099/core")
		if err != nil {
			issue = err
			return
		}

		// set synchronous session
		ses.SetMode(mgo.Monotonic, true)

		// store the session to the singleton instance
		instance.session = ses
	})

	return instance.session, issue
}

// Connect gh
func Connect(done chan error, sig chan struct{}) {
	fmt.Println("Connecting to database") // use logger for printing..
	session, err := GetInstance()
	if err != nil {
		fmt.Println(fmt.Errorf("DATABASE-ERROR: " + err.Error()))
		done <- err
		return
	}

	// the run loop
	go func() {
		for {
			select {
			case <-sig: // stop service loop
				fmt.Println("\tDatabase session .. ")
				session.Close()

				done <- nil

				return

			default:
			}
		}
	}()
}
