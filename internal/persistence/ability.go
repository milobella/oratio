package persistence

import (
	"context"
	"time"

	"github.com/milobella/oratio/internal/config"
	"github.com/milobella/oratio/internal/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AbilityDAO interface {
	CreateOrUpdate(ability *model.Ability) (*model.Ability, error)
	GetAll() ([]*model.Ability, error)
	GetByIntent(intent string) ([]*model.Ability, error)
}

type abilityDAOMongo struct {
	client     *mongo.Client
	url        string
	database   string
	collection string
	timeout    time.Duration
}

func NewAbilityDAOMongo(conf config.DatabaseConfig, timeout time.Duration) (AbilityDAO, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.MongoUrl))
	return &abilityDAOMongo{
		client:     client,
		url:        conf.MongoUrl,
		database:   conf.MongoDatabase,
		collection: conf.MongoCollection,
		timeout:    timeout,
	}, err
}

func (dao *abilityDAOMongo) CreateOrUpdate(ability *model.Ability) (*model.Ability, error) {
	collection := dao.client.Database(dao.database).Collection(dao.collection)

	opts := options.FindOneAndReplace().SetUpsert(true)
	ctx, _ := context.WithTimeout(context.Background(), dao.timeout)
	filter := bson.D{{"name", ability.Name}}

	result := collection.FindOneAndReplace(ctx, filter, ability, opts)

	var foundAbility *model.Ability
	err := result.Decode(foundAbility)
	return foundAbility, err
}

func (dao *abilityDAOMongo) GetAll() ([]*model.Ability, error) {
	collection := dao.client.Database(dao.database).Collection(dao.collection)
	ctx, _ := context.WithTimeout(context.Background(), dao.timeout)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		dao.logError(err, "Error creating the database cursor")
		return []*model.Ability{}, err
	}
	var results []*model.Ability
	if err = cursor.All(ctx, &results); err != nil {
		dao.logError(err, "Error getting results from the cursor")
		return []*model.Ability{}, err
	}
	return results, nil
}

func (dao *abilityDAOMongo) GetByIntent(intent string) ([]*model.Ability, error) {
	collection := dao.client.Database(dao.database).Collection(dao.collection)
	ctx, _ := context.WithTimeout(context.Background(), dao.timeout)
	cursor, err := collection.Find(ctx, bson.M{"intents": intent})
	if err != nil {
		dao.logError(err, "Error creating the database cursor")
		return []*model.Ability{}, err
	}
	var results []*model.Ability
	if err = cursor.All(ctx, &results); err != nil {
		dao.logError(err, "Error getting results from the cursor")
		return []*model.Ability{}, err
	}
	return results, nil
}

func (dao *abilityDAOMongo) logError(err error, message string) {
	logrus.WithError(err).
		WithField("url", dao.url).
		WithField("database", dao.database).
		WithField("collection", dao.collection).
		Error(message)
}
