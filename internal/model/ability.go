package model

// Ability is used in request/response body of the /api/v1/abilities endpoint
type Ability struct {
	Name    string   `json:"name"`
	Host    string   `json:"host"`
	Port    int      `json:"port"`
	Intents []string `json:"intents"`
}

// Abilities is the response body of the /api/v1/abilities endpoint (when no particular "from" query param is selected)
type Abilities struct {
	Cache    []*Ability `json:"cache"`
	Database []*Ability `json:"database"`
	Config   []*Ability `json:"config"`
}
