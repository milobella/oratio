package ability

import (
	"fmt"

	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/internal/model"
	"github.com/milobella/oratio/pkg/ability"
	"github.com/milobella/oratio/pkg/cerebro"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// Used to compute an approximate size of the map that will welcome the clients (one client by ability and by intent)
const approximateIntentsByAbility = 3

type Service interface {
	RequestAbility(nlu cerebro.NLU, context ability.Context, device ability.Device) *ability.Response
	GetCacheAbilities() ([]*model.Ability, error)
	GetDatabaseAbilities() ([]*model.Ability, error)
	GetConfigAbilities() ([]*model.Ability, error)
	GetAllAbilities() (*model.Abilities, error)
	CreateOrUpdate(ability *model.Ability) (*model.Ability, error)
}

// clients is used to store and index clients computed from abilities. It is used only for abilities coming
// from configuration because cache and database have their own indexation.
// Moreover, we don't want to bump all clients in the memory. We build clients from database data in a lazy mode.
// Configuration is just here in a last resort, if database is not accessible for example.
type clients = map[string]*ability.Client

func newClients(configAbilities []model.Ability) clients {
	clientsMap := make(map[string]*ability.Client, len(configAbilities)*(approximateIntentsByAbility+1))
	for _, ab := range configAbilities {
		client := ability.NewClient(ab.Host, ab.Port, ab.Name)
		for _, intent := range ab.Intents {
			clientsMap[intent] = client
		}
		clientsMap[client.Name] = client
	}
	return clientsMap
}

type serviceImpl struct {
	dao               DAO
	clientsCache      *cache.Cache
	clientsFromConfig clients
	stopIntent        string
}

func NewService(dao DAO, conf config.Abilities) Service {
	return &serviceImpl{
		dao:               dao,
		clientsCache:      cache.New(conf.Cache.Expiration, conf.Cache.CleanupInterval),
		clientsFromConfig: newClients(conf.List),
		stopIntent:        conf.StopIntent,
	}
}

// getBestIntentOrAbility computes the best intent or ability from the nlu but also from the context.
// Indeed, if the context says that a slot filling mechanism is in place, we force the future request to be sent to
// last ability. Except if the best intent is the stop intent, in which case we return the stop intent and everything
// will be stopped.
func (s *serviceImpl) getBestIntentOrAbility(nlu cerebro.NLU, ctx ability.Context) string {
	var bestIntent string
	if len(nlu.Intents) != 0 && nlu.BestIntent != "" {
		bestIntent = nlu.BestIntent
	}

	if ctx.SlotFilling == nil || bestIntent == s.stopIntent {
		return bestIntent
	}

	return ctx.LastAbility
}

// RequestAbility Call ability corresponding to the intent resolved by cerebro.
func (s *serviceImpl) RequestAbility(nlu cerebro.NLU, ctx ability.Context, device ability.Device) *ability.Response {

	intentOrAbility := s.getBestIntentOrAbility(nlu, ctx)

	// TODO put personal request in anima
	if intentOrAbility == "HELLO" {
		return ability.NewSimpleResponse("Hello")
	}

	if intentOrAbility == s.stopIntent {
		return ability.NewSimpleResponse("")
	}

	if client, ok := s.resolveClient(intentOrAbility); ok {
		if response, err := client.CallAbility(ability.Request{Nlu: nlu, Context: ctx, Device: device}); err == nil {
			if err = s.clientsCache.Add(intentOrAbility, client, cache.DefaultExpiration); err != nil {
				logrus.
					WithError(err).
					WithField("intentOrAbility", intentOrAbility).
					WithField("client", client.Name).
					Warning("An error occurred on adding the client in the cache.")
			}
			response.Context.LastAbility = client.Name
			return response
		}
	}

	return ability.NewSimpleResponse("I didn't find any ability corresponding to your request.")
}

// GetCacheAbilities fetch the abilities from the cache.
func (s *serviceImpl) GetCacheAbilities() ([]*model.Ability, error) {
	cacheItems := s.clientsCache.Items()
	abilities := make([]*model.Ability, len(cacheItems))
	for intent, item := range cacheItems {
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
	return abilities, nil
}

// GetDatabaseAbilities fetch the abilities from the database.
func (s *serviceImpl) GetDatabaseAbilities() ([]*model.Ability, error) {
	return s.dao.GetAll()
}

// GetConfigAbilities fetch the abilities from the configuration.
func (s *serviceImpl) GetConfigAbilities() ([]*model.Ability, error) {
	abilities := make([]*model.Ability, len(s.clientsFromConfig))
	for intent, client := range s.clientsFromConfig {
		abilities = append(abilities, &model.Ability{
			Name:    client.Name,
			Host:    client.Host,
			Port:    client.Port,
			Intents: []string{intent},
		})
	}
	return abilities, nil
}

// GetAllAbilities fetch the abilities from the every place (cache, database, config).
func (s *serviceImpl) GetAllAbilities() (*model.Abilities, error) {
	result := &model.Abilities{}
	var err error
	result.Cache, err = s.GetCacheAbilities()
	if err != nil {
		logrus.WithError(err).Error("An error occurred while fetching Abilities from cache")
		return nil, err
	}
	result.Database, err = s.GetDatabaseAbilities()
	if err != nil {
		logrus.WithError(err).Error("An error occurred while fetching Abilities from database")
		return nil, err
	}
	result.Config, err = s.GetConfigAbilities()
	if err != nil {
		logrus.WithError(err).Error("An error occurred while fetching Abilities from config")
		return nil, err
	}
	return result, nil
}
func (s *serviceImpl) CreateOrUpdate(ability *model.Ability) (*model.Ability, error) {
	return s.dao.CreateOrUpdate(ability)
}

func (s *serviceImpl) resolveClient(intentOrAbility string) (*ability.Client, bool) {
	// Resolve from cache
	if cachedClient, ok := s.clientsCache.Get(intentOrAbility); ok {
		client := cachedClient.(*ability.Client)
		logResolvedClientFrom("cache", intentOrAbility, client.Name)
		return client, true
	}

	// If not found, resolve from database
	abilities, err := s.dao.GetByIntent(intentOrAbility)
	if err == nil && len(abilities) > 0 {
		client := ability.NewClient(abilities[0].Host, abilities[0].Port, abilities[0].Name)
		logResolvedClientFrom("database", intentOrAbility, client.Name)
		return client, true
	}

	// If not found, resolve from config
	if client, ok := s.clientsFromConfig[intentOrAbility]; ok {
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
