package anima

type NLG struct {
	Sentence	string  	`json:"sentence,omitempty"`
	Params      map[string]string	`json:"params,omitempty"`
}