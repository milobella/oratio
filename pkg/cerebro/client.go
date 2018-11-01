package cerebro

type Client struct {

}

// TODO call effectively cerebro
func (c Client) UnderstandText(t string) NLU {
	if t == "Bonjour" {
		return NLU {"hello", "hello", []string{""}}
	}
	return NLU {"error", "error", []string{""}}
}