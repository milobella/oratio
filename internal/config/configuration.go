package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// Configuration
type Configuration struct {
	Server            ServerConfiguration
	Cerebro           CerebroConfiguration
	Anima             AnimaConfiguration
	Abilities         []AbilityConfiguration
	AbilitiesDatabase AbilitiesDatabaseConfiguration `mapstructure:"abilities_database"`
	AppSecret         string                         `mapstructure:"app_secret"`
}

// fun String() : Serialization function of Configuration (for logging)
func (c Configuration) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Configuration serialization error %s", err)
	}
	return string(b)
}

type ServerConfiguration struct {
	Port     int
	LogLevel string `mapstructure:"log_level"`
}

type CerebroConfiguration struct {
	Host string
	Port int
}

type AnimaConfiguration struct {
	Host string
	Port int
}

type AbilitiesDatabaseConfiguration struct {
	MongoUrl string `mapstructure:"mongo_url"`
}

type AbilityConfiguration struct {
	Name    string
	Host    string
	Port    int
	Intents []string
}
