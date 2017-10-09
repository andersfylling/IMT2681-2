package currencyfetcher

// Info contains the response from http://api.fixer.io/latest?base=EUR
// with json values
type Info struct {
	Base  string         `json:"base"`
	Date  string         `json:"date"`
	Rates map[string]int `json:"rates"`
}

func (i *Info) NewService() *Info {
	return &Info{} // TODO
}

// implement ServiceInterface methods
func (i *Info) load() {}

func (i *Info) run() {}
