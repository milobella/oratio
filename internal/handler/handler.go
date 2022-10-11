package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/milobella/oratio/internal/ability"
	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"github.com/sirupsen/logrus"
)

// New initiates all the handlers and their dependencies from the config.
// It will return an object containing directly the HandlerFunc to register to the server.
func New(conf *config.Config) *Handler {
	// Initialize clients
	cerebroClient := cerebro.NewClient(conf.Cerebro.Host, conf.Cerebro.Port, conf.Cerebro.UnderstandEndpoint)
	animaClient := anima.NewClient(conf.Anima.Host, conf.Anima.Port, conf.Anima.RestituteEndpoint)

	// Build the ability service. It will manage the DB and request the different abilities.
	// TODO(Célian): The DAO being not initialized shouldn't be a reason of not being up. We should have a retry mechanism + health endpoint
	abilityDAO, err := ability.NewMongoDAO(conf.Abilities.Database, 3*time.Second)
	if err != nil {
		logrus.WithError(err).Fatalf("Error initializing the Ability DAO.")
	}
	abilityService := ability.NewService(abilityDAO, conf.Abilities)

	// Build the handlers
	abilityHandler := NewAbility(abilityService)
	textHandler := NewText(cerebroClient, animaClient, abilityService)

	return &Handler{
		Text:          textHandler.Send,
		GetAbilities:  abilityHandler.Get,
		CreateAbility: abilityHandler.Create,
	}
}

type Handler struct {
	Text          echo.HandlerFunc
	GetAbilities  echo.HandlerFunc
	CreateAbility echo.HandlerFunc
}
