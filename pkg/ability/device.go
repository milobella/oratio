package ability

// Device information
type Device struct {
	// Some dynamic information sent with each request
	State       map[string]interface{} `json:"state,omitempty"`
	Instruments []interface{}          `json:"instruments,omitempty"`
}
