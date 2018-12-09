package cerebro

type NLU struct {
	BestIntent string
	Intents    []Intent
	Entities   []Entity
}

type Intent struct {
	Label string
	Score float32
}

type Entity struct {
	Label string
	Text  string
}
