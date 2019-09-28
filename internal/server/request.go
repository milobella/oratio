package server

import "milobella.com/gitlab/milobella/oratio/pkg/ability"

// RequestBody is the main request body or oratio
type RequestBody struct {
	Text    string          `json:"text,omitempty"`
	Context ability.Context `json:"context,omitempty"`
}
