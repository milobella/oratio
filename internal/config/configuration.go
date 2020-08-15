package config

import (
	"encoding/json"
	"github.com/milobella/oratio/internal/models"
	"github.com/sirupsen/logrus"
	"time"
)

// Configuration
type Configuration struct {
	Server            ServerConfiguration
	Cerebro           CerebroConfiguration
	Anima             AnimaConfiguration
	Abilities         []models.Ability
	AbilitiesCache    AbilitiesCacheConfiguration    `mapstructure:"abilities_cache"`
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
	UnderstandEndpoint string `mapstructure:"understand_endpoint"`
}

type AnimaConfiguration struct {
	Host string
	Port int
	RestituteEndpoint string `mapstructure:"restitute_endpoint"`
}

type AbilitiesDatabaseConfiguration struct {
	MongoUrl string `mapstructure:"mongo_url"`
}
type AbilitiesCacheConfiguration struct {
	Expiration      time.Duration
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}
