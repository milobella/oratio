package service

import (
	"fmt"

	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/internal/model"
	"github.com/milobella/oratio/internal/persistence"
	"github.com/milobella/oratio/pkg/ability"
	"github.com/milobella/oratio/pkg/cerebro"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// Used to compute an approximative size of the map that will welcome the clients (one client by ability and by intent)
const approximativeIntentsByAbility = 3

type Abilities struct {
	Cache    []*model.Ability `json:"cache"`
	Database []*model.Ability `json:"database"`
	Config   []*model.Ability `json:"config"`
}

type AbilityService interface {
	RequestAbility(nlu cerebro.NLU, context ability.Context, device ability.Device) *ability.Response
	GetCacheAbilities() ([]*model.Ability, error)
	GetDatabaseAbilities() ([]*model.Ability, error)
	GetConfigAbilities() ([]*model.Ability, error)
	GetAllAbilities() (*Abilities, error)
}

// abilityClients is used to store and index clients computed from abilities. It is used only for abilities coming
// from configuration because cache and database have their own indexation.
// Moreover, we don't want to bump all clients in the memory. We build clients from database data in a lazy mode.
// Configuration is just here in a last resort, if database is not accessible for example.
type abilityClients = map[string]*ability.Client

func newAbilityClients(configAbilities []model.Ability) abilityClients {
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
	dao                      persistence.AbilityDAO
	clientsCache             *cache.Cache
	abilityClientsFromConfig abilityClients
	stopIntent               string
}

func NewAbilityService(dao persistence.AbilityDAO, abilitiesConfig config.AbilitiesConfig) AbilityService {
	return &abilityServiceImpl{
		dao:                      dao,
		clientsCache:             cache.New(abilitiesConfig.Cache.Expiration, abilitiesConfig.Cache.CleanupInterval),
		abilityClientsFromConfig: newAbilityClients(abilitiesConfig.List),
		stopIntent:               abilitiesConfig.StopIntent,
	}
}

// getBestIntentOrAbility computes the best intent or ability from the nlu but also from the context.
// Indeed, if the context says that a slot filling mechanism is in place, we force the future request to be sent to
// last ability. Except if the best intent is the stop intent, in which case we return the stop intent and everything
// will be stopped.
func (a *abilityServiceImpl) getBestIntentOrAbility(nlu cerebro.NLU, ctx ability.Context) string {
	var bestIntent string
	if len(nlu.Intents) != 0 && nlu.BestIntent != "" {
		bestIntent = nlu.BestIntent
	}

	if ctx.SlotFilling == nil || bestIntent == a.stopIntent {
		return bestIntent
	}

	return ctx.LastAbility
}

// RequestAbility Call ability corresponding to the intent resolved by cerebro.
func (a *abilityServiceImpl) RequestAbility(nlu cerebro.NLU, ctx ability.Context, device ability.Device) *ability.Response {

	intentOrAbility := a.getBestIntentOrAbility(nlu, ctx)

	// TODO put personal request in anima
	if intentOrAbility == "HELLO" {
		return ability.NewSimpleResponse("Hello")
	}

	if intentOrAbility == a.stopIntent {
		return ability.NewSimpleResponse("")
	}

	if client, ok := a.resolveClient(intentOrAbility); ok {
		if response, err := client.CallAbility(ability.Request{Nlu: nlu, Context: ctx, Device: device}); err == nil {
			if err = a.clientsCache.Add(intentOrAbility, client, cache.DefaultExpiration); err != nil {
				logrus.
					WithError(err).
					WithField("intentOrAbility", intentOrAbility).
					WithField("client", client.Name).
					Warning("An error occurred on adding the client in the cache.")
			} else {
				response.Context.LastAbility = client.Name
			}
			return response
		}
	}

	return ability.NewSimpleResponse("I didn't find any ability corresponding to your request.")
}

// GetCacheAbilities fetch the abilities from the cache.
func (a *abilityServiceImpl) GetCacheAbilities() ([]*model.Ability, error) {
	var abilities []*model.Ability
	for intent, item := range a.clientsCache.Items() {
		client, ok := item.Object.(*ability.Client)
		if !ok {
			return nil, fmt.Errorf("error casting cache entry into %T", &ability.Client{})
		}
		abilities = append(abilities, &model.Ability{
			Name:    client.Name,
			Host:    client.Host,
			Port:    client.Port,
			Intents: []string{intent},
		})
	}
	if abilities == nil {
		return []*model.Ability{}, nil
	} else {
		return abilities, nil
	}
}

// GetDatabaseAbilities fetch the abilities from the database.
func (a *abilityServiceImpl) GetDatabaseAbilities() ([]*model.Ability, error) {
	return a.dao.GetAll()
}

// GetConfigAbilities fetch the abilities from the configuration.
func (a *abilityServiceImpl) GetConfigAbilities() ([]*model.Ability, error) {
	var abilities []*model.Ability
	for intent, client := range a.abilityClientsFromConfig {
		abilities = append(abilities, &model.Ability{
			Name:    client.Name,
			Host:    client.Host,
			Port:    client.Port,
			Intents: []string{intent},
		})
	}
	if abilities == nil {
		return []*model.Ability{}, nil
	} else {
		return abilities, nil
	}
}

// GetAllAbilities fetch the abilities from the every place (cache, database, config).
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
