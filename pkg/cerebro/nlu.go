package cerebro

type NLU struct {
	Intent    string   `json:"action,omitempty"`
	Parameter []string `json:"parameter,omitempty"`
}