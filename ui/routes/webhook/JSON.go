package webhook

// CreateJSON For recieving request create a listener to currency changes
type CreateJSON struct {
	URL    string  `json:"webhookURL"`
	Base   string  `json:"baseCurrency"`
	Target string  `json:"targetCurrency"`
	Min    float64 `json:"minTriggerValue"`
	Max    float64 `json:"maxTriggerValue"`
}

// InvokeJSON For invoking webhooks
type InvokeJSON struct {
	Base   string  `json:"baseCurrency"`
	Target string  `json:"targetCurrency"`
	Min    float64 `json:"minTriggerValue"`
	Max    float64 `json:"maxTriggerValue"`
}
