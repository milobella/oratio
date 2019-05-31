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

func (nlu *NLU) GetBestIntentOr(fallback string) string {
	if len(nlu.Intents) != 0 && nlu.BestIntent != ""  {
		return nlu.BestIntent
	}

	return fallback
}
