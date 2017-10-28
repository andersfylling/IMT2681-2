package currency

// Collection to be stored in database
type Collection struct {
	Base  string         `json:"base"`
	Date  string         `json:"date"`
	Rates map[string]int `json:"rates"`
}

func Save() {
	// Find
	// if matched, Update
	// if not found, use Create

}

func Create() {}
func Update() {}
func Delete() {}
func Find()   {}
