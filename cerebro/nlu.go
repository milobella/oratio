package cerebro

type NLU struct {
	Skill		string  	`json:"skill,omitempty"`
	Action 		string 		`json:"action,omitempty"`
	Parameter 	[]string 	`json:"parameter,omitempty"`
}