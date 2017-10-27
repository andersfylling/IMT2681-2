package latest

// PostJSON For recieving request create a listener to currency changes
// calculate the avg for the last seven days
type PostJSON struct {
	Base   string `json:"baseCurrency"`
	Target string `json:"targetCurrency"`
}
