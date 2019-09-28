package server

import (
	"milobella.com/gitlab/milobella/oratio/pkg/ability"
	"milobella.com/gitlab/milobella/oratio/pkg/anima"
	"milobella.com/gitlab/milobella/oratio/pkg/cerebro"
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
