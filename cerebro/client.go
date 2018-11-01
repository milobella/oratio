package cerebro

type CerebroClient struct {

}

func (c CerebroClient) UnderstandText(t string) NLU {
	n := NLU {"clock", "get_time", []string{""}}
	return n
}