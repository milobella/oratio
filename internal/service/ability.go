package service

import (
	"fmt"
	"github.com/milobella/oratio/internal/models"
	"github.com/milobella/oratio/pkg/ability"
	"github.com/milobella/oratio/pkg/anima"
	"github.com/milobella/oratio/pkg/cerebro"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
	"time"
)

// Used to compute an approximative size of the map that will welcome the clients (one client by ability and by intent)
const approximativeIntentsByAbility = 3

type Abilities struct {
	Cache    []*models.Ability `json:"cache"`
	Database []*models.Ability `json:"database"`
	Config   []*models.Ability `json:"config"`
}

type AbilityService interface {
	RequestAbility(nlu cerebro.NLU, context ability.Context, device ability.Device) ability.Response
	GetCacheAbilities() ([]*models.Ability, error)
	GetDatabaseAbilities() ([]*models.Ability, error)
	GetConfigAbilities() ([]*models.Ability, error)
	GetAllAbilities() (*Abilities, error)
}

// abilityClients is used to store and index clients computed from abilities. It is used only for abilities coming
// from configuration because cache and database have their own indexation.
// Moreover, we don't want to bump all clients in the memory. We build clients from database data in a lazy mode.
// Configuration is just here in a last resort, if database is not accessible for example.
type abilityClients = map[string]*ability.Client

func newAbilityClients(configAbilities []models.Ability) abilityClients {
	clientsMap := make(map[string]*ability.Client, len(configAbilities)*approximativeIntentsByAbility)
	for _, ab := range configAbilities {
		client := ability.NewClient(ab.Host, ab.Port, ab.Name)
		for _, intent := range ab.Intents {
			clientsMap[intent] = client
		}
	}
	return clientsMap
}

type abilityServiceImpl struct {
	dao                      models.AbilityDAO
	clientsCache             *cache.Cache
	abilityClientsFromConfig abilityClients
}

func NewAbilityService(dao models.AbilityDAO, configAbilities []models.Ability, defaultExpiration, cleanupInterval time.Duration) AbilityService {
	return &abilityServiceImpl{
		dao:                      dao,
		clientsCache:             cache.New(defaultExpiration, cleanupInterval),
		abilityClientsFromConfig: newAbilityClients(configAbilities),
	}
}

// RequestAbility: Call ability corresponding to the intent resolved by cerebro.
func (a *abilityServiceImpl) RequestAbility(nlu cerebro.NLU, context ability.Context, device ability.Device) ability.Response {

	intentOrAbility := nlu.GetBestIntentOr(context.LastAbility)

	// TODO put personal request in anima
	if intentOrAbility == "HELLO" {
		return ability.Response{
			Nlg:          anima.NLG{Sentence: "Hello"},
			Visu:         nil,
			AutoReprompt: false,
			Context:      ability.Context{},
		}
	}

	if client, ok := a.resolveClient(intentOrAbility); ok {
		if response, err := client.CallAbility(ability.Request{Nlu: nlu, Context: context, Device: device}); err == nil {
			if err = a.clientsCache.Add(intentOrAbility, client, cache.DefaultExpiration); err != nil {
				logrus.
					WithError(err).
					WithField("intentOrAbility", intentOrAbility).
					WithField("client", client.Name).
					Warning("An error occurred on adding the client in the cache.")
			}
			return *response
		}
	}

	return ability.Response{
		Nlg:          anima.NLG{Sentence: "Oups !"},
		Visu:         nil,
		AutoReprompt: false,
		Context:      ability.Context{},
	}
}

// GetCacheAbilities: Fetch the abilities from the cache.
func (a *abilityServiceImpl) GetCacheAbilities() ([]*models.Ability, error) {
	var abilities []*models.Ability
	for intent, item := range a.clientsCache.Items() {
		client, ok := item.Object.(*ability.Client)
		if !ok {
			return nil, fmt.Errorf("error casting cache entry into %T", &ability.Client{})
		}
		abilities = append(abilities, &models.Ability{
			Name:    client.Name,
			Host:    client.Host,
			Port:    client.Port,
			Intents: []string{intent},
		})
	}
	if abilities == nil {
		return []*models.Ability{}, nil
	} else {
		return abilities, nil
	}
}

// GetDatabaseAbilities: Fetch the abilities from the database.
func (a *abilityServiceImpl) GetDatabaseAbilities() ([]*models.Ability, error) {
	return a.dao.GetAll()
}

// GetConfigAbilities: Fetch the abilities from the configuration.
func (a *abilityServiceImpl) GetConfigAbilities() ([]*models.Ability, error) {
	var abilities []*models.Ability
	for intent, client := range a.abilityClientsFromConfig {
		abilities = append(abilities, &models.Ability{
			Name:    client.Name,
			Host:    client.Host,
			Port:    client.Port,
			Intents: []string{intent},
		})
	}
	if abilities == nil {
		return []*models.Ability{}, nil
	} else {
		return abilities, nil
	}
}

// GetAllAbilities: Fetch the abilities from the every places (cache, database, config).
func (a *abilityServiceImpl) GetAllAbilities() (*Abilities, error) {
	cacheAbilities, err := a.GetCacheAbilities()
	if err != nil {
		logrus.WithError(err).Error("An error occurred while fetching Abilities from cache")
		return nil, err
	}
	databaseAbilities, err := a.GetDatabaseAbilities()
	if err != nil {
		logrus.WithError(err).Error("An error occurred while fetching Abilities from database")
		return nil, err
	}
	configAbilities, err := a.GetConfigAbilities()
	if err != nil {
		logrus.WithError(err).Error("An error occurred while fetching Abilities from config")
		return nil, err
	}
	return &Abilities{
		Cache:    cacheAbilities,
		Database: databaseAbilities,
		Config:   configAbilities,
	}, nil
}

func (a *abilityServiceImpl) resolveClient(intentOrAbility string) (*ability.Client, bool) {
	// Resolve from cache
	if cachedClient, ok := a.clientsCache.Get(intentOrAbility); ok {
		client := cachedClient.(*ability.Client)
		logResolvedClientFrom("cache", intentOrAbility, client.Name)
		return client, true
	}

	// If not found, resolve from database
	abilities, err := a.dao.GetByIntent(intentOrAbility)
	if err == nil && len(abilities) > 0 {
		client := ability.NewClient(abilities[0].Host, abilities[0].Port, abilities[0].Name)
		logResolvedClientFrom("database", intentOrAbility, client.Name)
		return client, true
	}

	// If not found, resolve from config
	if client, ok := a.abilityClientsFromConfig[intentOrAbility]; ok {
		logResolvedClientFrom("configuration", intentOrAbility, client.Name)
		return client, true
	}

	logrus.
		WithError(err).
		WithField("intentOrAbility", intentOrAbility).
		Error("Didn't find any ability for this intent or ability name.")
	return nil, false
}

func logResolvedClientFrom(location string, intentOrAbility string, client string) {
	logrus.
		WithField("intentOrAbility", intentOrAbility).
		WithField("client", client).
		Debugf("Resolved the client from %s to request ability.", location)
}
