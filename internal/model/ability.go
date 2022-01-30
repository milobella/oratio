package model

type Ability struct {
	Name    string   `json:"name"`
	Host    string   `json:"host"`
	Port    int      `json:"port"`
	Intents []string `json:"intents"`
}
