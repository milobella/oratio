package server

import (
	"github.com/milobella/oratio/pkg/ability"
)

// RequestBody is the main request body
type RequestBody struct {
	Text    string          `json:"text,omitempty"`
	Context ability.Context `json:"context,omitempty"`
	Device  ability.Device  `json:"device,omitempty"`
}
