package service

// Interface holds info for runnable services
// TODO-if-needed: let services say if they can run in parallel or need to be
//                 run after a specific service or services
type Interface interface {
	NewService() Interface
	Title() string    // I really wish interfaces supported const members
	Load()            // Load needed data into memory
	Empty()           // Empty loaded data from memory
	Run()             // Run the service action
	Timeout() float64 // Time left before the service can run again.
	// TODO: have a way to store data ?
}
