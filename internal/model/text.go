package model

import "github.com/milobella/oratio/pkg/ability"

// TextRequest is the request body of /api/v1/talke/text endpoint
type TextRequest struct {
	Text    string          `json:"text,omitempty"`
	Context ability.Context `json:"context,omitempty"`
	Device  ability.Device  `json:"device,omitempty"`
}

// TextResponse is the response body of /api/v1/talke/text endpoint
type TextResponse struct {
	Vocal        string      `json:"vocal,omitempty"`
	Visu         interface{} `json:"visu,omitempty"`
	Actions      interface{} `json:"actions,omitempty"`
	AutoReprompt bool        `json:"auto_reprompt,omitempty"`
	Context      interface{} `json:"context,omitempty"`
}
