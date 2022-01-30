package config

import (
	"encoding/json"
	"github.com/milobella/oratio/internal/model"
	"github.com/sirupsen/logrus"
	"time"
)

type Configuration struct {
	Server            ServerConfiguration
	Tracing           TracingConfiguration
	Auth              AuthConfiguration
	Cerebro           CerebroConfiguration
	Anima             AnimaConfiguration
	Abilities         []model.Ability
	AbilitiesCache    AbilitiesCacheConfiguration    `mapstructure:"abilities_cache"`
	AbilitiesDatabase AbilitiesDatabaseConfiguration `mapstructure:"abilities_database"`
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
	ServiceName string `mapstructure:"service_name"`
	Port        int
	LogLevel    logrus.Level `mapstructure:"log_level"`
}

type TracingConfiguration struct {
	ServiceName         string `mapstructure:"service_name"`
	JaegerAgentHostName string `mapstructure:"jaeger_agent_hostname"`
	JaegerAgentPort     int    `mapstructure:"jaeger_agent_port"`
}

type AuthConfiguration struct {
	AppSecret string `mapstructure:"app_secret"`
}

type CerebroConfiguration struct {
	Host               string
	Port               int
	UnderstandEndpoint string `mapstructure:"understand_endpoint"`
}

type AnimaConfiguration struct {
	Host              string
	Port              int
	RestituteEndpoint string `mapstructure:"restitute_endpoint"`
}

type AbilitiesDatabaseConfiguration struct {
	MongoDatabase   string `mapstructure:"mongo_database"`
	MongoUrl        string `mapstructure:"mongo_url"`
	MongoCollection string `mapstructure:"mongo_collection"`
}
type AbilitiesCacheConfiguration struct {
	Expiration      time.Duration
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}
