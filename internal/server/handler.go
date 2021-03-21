package server

import (
	"github.com/labstack/echo"
	"github.com/milobella/oratio/internal/models"
	"github.com/milobella/oratio/internal/service"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"net/http"
)

type TextRequestHandler struct {
	CerebroClient  *cerebro.Client
	AnimaClient    *anima.Client
	AbilityService service.AbilityService
}

func (rh *TextRequestHandler) HandleTextRequest(c echo.Context) (err error) {
	// Read the request
	requestBody := new(RequestBody)
	if err = c.Bind(requestBody); err != nil {
		return
	}

	// Execute the processing flow
	nlu := rh.CerebroClient.UnderstandText(requestBody.Text)
	response := rh.AbilityService.RequestAbility(nlu, requestBody.Context, requestBody.Device)
	vocal := rh.AnimaClient.GenerateSentence(response.Nlg)

	// Build the response's body
	responseBody := &ResponseBody{Vocal: vocal, Visu: response.Visu, AutoReprompt: response.AutoReprompt, Context: response.Context}

	// Write it on the http response
	return c.JSON(http.StatusOK, responseBody)
}

type AbilityRequestHandler struct {
	AbilitDAO      models.AbilityDAO
	AbilityService service.AbilityService
}

func (rh *AbilityRequestHandler) HandleGetAllAbilityRequest(c echo.Context) (err error) {
	from := c.QueryParam("from")
	switch from {
	case "cache":
		if result, err := rh.AbilityService.GetCacheAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	case "database":
		if result, err := rh.AbilityService.GetDatabaseAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	case "config":
		if result, err := rh.AbilityService.GetConfigAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	default:
		if result, err := rh.AbilityService.GetAllAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	}
}

func (rh *AbilityRequestHandler) HandleCreateAbilityRequest(c echo.Context) (err error) {
	ability := new(models.Ability)
	if err = c.Bind(ability); err != nil {
		return
	}

	if result, err := rh.AbilitDAO.CreateOrUpdate(ability); err != nil {
		return echo.NewHTTPError(500, err.Error())
	} else {
		return c.JSON(http.StatusOK, result)
	}
}
