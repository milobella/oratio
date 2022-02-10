package ability

import "github.com/milobella/oratio/pkg/anima"

type Response struct {
	Nlg          anima.NLG   `json:"nlg,omitempty"`
	Visu         interface{} `json:"visu,omitempty"`
	Actions      interface{} `json:"actions,omitempty"`
	AutoReprompt bool        `json:"auto_reprompt,omitempty"`
	Context      Context     `json:"context,omitempty"`
}

func NewSimpleResponse(text string) *Response {
	return &Response{Nlg: anima.NLG{Sentence: text}}
}
