package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/milobella/oratio/internal/ability"
	"github.com/milobella/oratio/internal/model"
	"github.com/sirupsen/logrus"
)

func NewAbility(service ability.Service) Ability {
	return &abilityImpl{service: service}
}

type Ability interface {
	Get(c echo.Context) (err error)
	Create(c echo.Context) (err error)
}

type abilityImpl struct {
	service ability.Service
}

func (a *abilityImpl) Get(c echo.Context) (err error) {
	from := c.QueryParam("from")
	defer func() {
		if err != nil {
			logrus.WithError(err).WithField("from", from).Error("An error occurred while getting Abilities")
		}
	}()
	switch from {
	case "cache":
		if result, err := a.service.GetCacheAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	case "database":
		if result, err := a.service.GetDatabaseAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	case "config":
		if result, err := a.service.GetConfigAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	default:
		if result, err := a.service.GetAllAbilities(); err != nil {
			return echo.NewHTTPError(500, err.Error())
		} else {
			return c.JSON(http.StatusOK, result)
		}
	}
}

func (a *abilityImpl) Create(c echo.Context) error {
	futureAbility := new(model.Ability)
	if err := c.Bind(futureAbility); err != nil {
		return err
	}

	if result, err := a.service.CreateOrUpdate(futureAbility); err != nil {
		return echo.NewHTTPError(500, err.Error())
	} else {
		return c.JSON(http.StatusOK, result)
	}
}
