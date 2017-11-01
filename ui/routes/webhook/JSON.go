package webhook

// Invoke For invoking webhooks
type Invoke struct {
	Base   string  `json:"baseCurrency"`
	Target string  `json:"targetCurrency"`
	Min    float64 `json:"minTriggerValue"`
	Max    float64 `json:"maxTriggerValue"`
}
