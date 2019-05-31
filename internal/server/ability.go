package server

import (
	"fmt"
	"gitlab.milobella.com/milobella/oratio/pkg/ability"
	"gitlab.milobella.com/milobella/oratio/pkg/anima"
	"gitlab.milobella.com/milobella/oratio/pkg/cerebro"
	"time"
)

type AbilityService struct {
	Clients map[string]*ability.Client
}

// RequestAbility: Call ability corresponding to the intent resolved by cerebro.
func (acr *AbilityService) RequestAbility(nlu cerebro.NLU, context ability.Context) (anima.NLG, interface{}, bool, ability.Context) {

	intentOrAbility := nlu.GetBestIntentOr(context.LastAbility)

	// TODO put personal request in anima
	if intentOrAbility == "HELLO" {
		return anima.NLG{Sentence: "Hello"}, nil, false, ability.Context{}
	}

	// TODO put time request in clock ability
	if intentOrAbility == "GET_TIME" {
		now := time.Now()
		timeVal := fmt.Sprintf("%d h %d", now.Hour(), now.Minute())
		return anima.NLG{
			Sentence: "It is {{time}}",
			Params: []anima.NLGParam{{
				Name:  "time",
				Value: timeVal,
				Type:  "time",
			}}}, nil, false, ability.Context{}
	}

	if client, ok := acr.Clients[intentOrAbility]; ok {
		return doRequest(client, ability.Request{Nlu: nlu, Context: context})
	}

	return anima.NLG{Sentence: "Oups !"}, nil, false, ability.Context{}
}

func doRequest(client *ability.Client, request ability.Request) (anima.NLG, interface{}, bool, ability.Context) {
	nlg, visu, autoReprompt, context := client.CallAbility(request)
	context.LastAbility = client.Name
	return nlg, visu, autoReprompt, context
}
