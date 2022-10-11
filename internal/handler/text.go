package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/milobella/oratio/internal/ability"
	"github.com/milobella/oratio/internal/model"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
)

func NewText(cerebroClient *cerebro.Client, animaClient *anima.Client, abilityService ability.Service) Text {
	return &textImpl{
		CerebroClient:  cerebroClient,
		AnimaClient:    animaClient,
		AbilityService: abilityService,
	}
}

type Text interface {
	Send(c echo.Context) (err error)
}

type textImpl struct {
	CerebroClient  *cerebro.Client
	AnimaClient    *anima.Client
	AbilityService ability.Service
}

func (rh *textImpl) Send(c echo.Context) (err error) {
	// Read the request
	requestBody := new(model.TextRequest)
	if err = c.Bind(requestBody); err != nil {
		return
	}

	// Execute the processing flow
	nlu := rh.CerebroClient.UnderstandText(requestBody.Text)
	response := rh.AbilityService.RequestAbility(nlu, requestBody.Context, requestBody.Device)
	vocal := rh.AnimaClient.GenerateSentence(response.Nlg)

	// Build the response's body
	responseBody := &model.TextResponse{
		Vocal:        vocal,
		Visu:         response.Visu,
		AutoReprompt: response.AutoReprompt,
		Context:      response.Context,
		Actions:      response.Actions,
	}

	// Write it on the http response
	return c.JSON(http.StatusOK, responseBody)
}
