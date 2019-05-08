package anima

type NLG struct {
	Sentence string     `json:"sentence,omitempty"`
	Params   []NLGParam `json:"params,omitempty"`
}

type NLGParam struct {
	Name  string      `json:"name,omitempty"`
	Value interface{} `json:"value,omitempty"`
	Type  string      `json:"type,omitempty"`
}
