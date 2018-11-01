package anima

type Client struct {

}

// TODO call effectively anima
func (c Client) GenerateSentence(nlg NLG) string {
	return nlg.Sentence
}