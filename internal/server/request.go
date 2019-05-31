package server

import "gitlab.milobella.com/milobella/oratio/pkg/ability"

// RequestBody is the main request body or oratio
type RequestBody struct {
	Text    string          `json:"text,omitempty"`
	Context ability.Context `json:"context,omitempty"`
}
