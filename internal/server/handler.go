package server

import (
	"encoding/json"
	"github.com/milobella/oratio/internal/models"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"io/ioutil"
	"net/http"
)

type TextRequestHandler struct {
	CerebroClient  *cerebro.Client
	AnimaClient    *anima.Client
	AbilityService AbilityService
}

func (rh *TextRequestHandler) HandleTextRequest(w http.ResponseWriter, r *http.Request) {
	// Read the request
	requestBody, err := rh.readTextRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// Execute the processing flow
	nlu := rh.CerebroClient.UnderstandText(requestBody.Text)
	nlg, visu, autoReprompt, context := rh.AbilityService.RequestAbility(nlu, requestBody.Context, requestBody.Device)
	vocal := rh.AnimaClient.GenerateSentence(nlg)

	// Build the response's body
	responseBody := ResponseBody{Vocal: vocal, Visu: visu, AutoReprompt: autoReprompt, Context: context}

	// Write it on the http response
	if err = json.NewEncoder(w).Encode(responseBody); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (rh *TextRequestHandler) readTextRequest(r *http.Request) (req RequestBody, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &req)
	return
}

type AbilityRequestHandler struct {
	AbilitDAO models.AbilityDAO
}

func (rh *AbilityRequestHandler) HandleAbilityRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		rh.handleCreateAbilityRequest(w, r)
		break
	case "GET":
		rh.handleGetAllAbilityRequest(w, r)
		break
	}
}

func (rh *AbilityRequestHandler) handleCreateAbilityRequest(w http.ResponseWriter, r *http.Request) {
	ability, err := rh.readAbilityRequest(r)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if ability == nil {
		http.Error(w, "You gave an empty ability", 400)
		return
	}
	result, err := rh.AbilitDAO.CreateOrUpdate(ability)
	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), 500)
	}

}

func (rh *AbilityRequestHandler) handleGetAllAbilityRequest(w http.ResponseWriter, r *http.Request) {
	result, err := rh.AbilitDAO.GetAll()
	if err = json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (rh *AbilityRequestHandler) readAbilityRequest(r *http.Request) (req *models.Ability, err error) {
	b, err := ioutil.ReadAll(r.Body)
	defer func() {
		_ = r.Body.Close()
	}()
	if err != nil {
		return
	}

	var result models.Ability
	err = json.Unmarshal(b, &result)
	return &result, err
}
