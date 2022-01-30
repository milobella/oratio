package server

// ResponseBody is the main response body
type ResponseBody struct {
	Vocal        string      `json:"vocal,omitempty"`
	Visu         interface{} `json:"visu,omitempty"`
	AutoReprompt bool        `json:"auto_reprompt,omitempty"`
	Context      interface{} `json:"context,omitempty"`
}
