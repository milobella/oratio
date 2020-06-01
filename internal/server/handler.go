package server

import (
	"encoding/json"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"io/ioutil"
	"net/http"
)

type TextRequestHandler struct {
	CerebroClient  *cerebro.Client
	AnimaClient    *anima.Client
	AbilityService *AbilityService
}

func (trh *TextRequestHandler) HandleTextRequest(w http.ResponseWriter, r *http.Request) {
	// Read the request
	requestBody, err := readRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Execute the processing flow
	nlu := trh.CerebroClient.UnderstandText(requestBody.Text)
	nlg, visu, autoReprompt, context := trh.AbilityService.RequestAbility(nlu, requestBody.Context, requestBody.Device)
	vocal := trh.AnimaClient.GenerateSentence(nlg)

	// Build the body of the response
	responseBody := ResponseBody{Vocal: vocal, Visu: visu, AutoReprompt: autoReprompt, Context: context}

	// Write it on the http response
	if err = json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func readRequest(r *http.Request) (req RequestBody, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &req)
	return
}
