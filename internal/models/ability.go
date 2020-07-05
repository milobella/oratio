package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Ability struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int `json:"port"`
	Intents []string `json:"intents"`
}

type AbilityDAO interface {
	CreateOrUpdate(ability *Ability) (*Ability, error)
	GetAll() ([]*Ability, error)
	GetByIntent(intent string) ([]*Ability, error)
}

const (
	mongoDBCollection = "abilities"
)

type abilityDAOMongo struct {
	client *mongo.Client
	database string
	timeout time.Duration
}

func NewAbilityDAOMongo(url string, database string, timeout time.Duration) (AbilityDAO, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(url))
	return &abilityDAOMongo{client: client, database: database, timeout: timeout}, err
}

func (dao *abilityDAOMongo) CreateOrUpdate(ability *Ability) (*Ability, error) {
	collection := dao.client.Database(dao.database).Collection(mongoDBCollection)

	opts := options.FindOneAndReplace().SetUpsert(true)
	ctx, _ := context.WithTimeout(context.Background(), dao.timeout)
	filter := bson.D{{"name", ability.Name}}

	result := collection.FindOneAndReplace(ctx, filter, ability, opts)

	var foundAbility *Ability
	err := result.Decode(foundAbility)
	return foundAbility, err
}

func (dao *abilityDAOMongo) GetAll() ([]*Ability, error) {
	collection := dao.client.Database(dao.database).Collection(mongoDBCollection)
	ctx, _ := context.WithTimeout(context.Background(), dao.timeout)
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return []*Ability{}, err
	}
	var results	[]*Ability
	if err = cursor.All(ctx, &results); err != nil {
		return []*Ability{}, err
	}
	return results, nil
}

func (dao *abilityDAOMongo) GetByIntent(intent string) ([]*Ability, error) {
	collection := dao.client.Database(dao.database).Collection(mongoDBCollection)
	ctx, _ := context.WithTimeout(context.Background(), dao.timeout)
	cursor, err := collection.Find(ctx, bson.M{"intents": intent})
	if err != nil {
		return []*Ability{}, err
	}
	var results	[]*Ability
	if err = cursor.All(ctx, &results); err != nil {
		return []*Ability{}, err
	}
	return results, nil
}
