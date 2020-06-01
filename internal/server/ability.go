package server

import (
	"github.com/milobella/oratio/pkg/ability"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
)

type AbilityService struct {
	Clients map[string]*ability.Client
}

// RequestAbility: Call ability corresponding to the intent resolved by cerebro.
func (acr *AbilityService) RequestAbility(nlu cerebro.NLU, context ability.Context, device ability.Device) (anima.NLG, interface{}, bool, ability.Context) {

	intentOrAbility := nlu.GetBestIntentOr(context.LastAbility)

	// TODO put personal request in anima
	if intentOrAbility == "HELLO" {
		return anima.NLG{Sentence: "Hello"}, nil, false, ability.Context{}
	}

	if client, ok := acr.Clients[intentOrAbility]; ok {
		return doRequest(client, ability.Request{Nlu: nlu, Context: context, Device: device})
	}

	return anima.NLG{Sentence: "Oups !"}, nil, false, ability.Context{}
}

func doRequest(client *ability.Client, request ability.Request) (anima.NLG, interface{}, bool, ability.Context) {
	nlg, visu, autoReprompt, context := client.CallAbility(request)
	context.LastAbility = client.Name
	return nlg, visu, autoReprompt, context
}
