package ability

import "github.com/milobella/oratio/pkg/cerebro"

type Request struct {
	Nlu     cerebro.NLU `json:"nlu,omitempty"`
	Context Context     `json:"context,omitempty"`
}


