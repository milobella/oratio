package config

import (
	"encoding/json"
	"time"

	"github.com/milobella/oratio/internal/model"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server    ServerConfig
	Tracing   TracingConfig
	Auth      AuthConfig
	Cerebro   CerebroConfig
	Anima     AnimaConfig
	Abilities AbilitiesConfig
}

// fun String() : Serialization function of Config (for logging)
func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Configuration serialization error %s", err)
	}
	return string(b)
}

type ServerConfig struct {
	ServiceName string `mapstructure:"service_name"`
	Port        int
	LogLevel    string `mapstructure:"log_level"`
}

type TracingConfig struct {
	ServiceName         string `mapstructure:"service_name"`
	JaegerAgentHostName string `mapstructure:"jaeger_agent_hostname"`
	JaegerAgentPort     int    `mapstructure:"jaeger_agent_port"`
}

type AuthConfig struct {
	AppSecret string `mapstructure:"app_secret"`
}

type CerebroConfig struct {
	Host               string
	Port               int
	UnderstandEndpoint string `mapstructure:"understand_endpoint"`
}

type AnimaConfig struct {
	Host              string
	Port              int
	RestituteEndpoint string `mapstructure:"restitute_endpoint"`
}

type AbilitiesConfig struct {
	List       []model.Ability
	Cache      CacheConfig
	Database   DatabaseConfig
	StopIntent string `mapstructure:"stop_intent"`
}

type DatabaseConfig struct {
	MongoDatabase   string `mapstructure:"mongo_database"`
	MongoUrl        string `mapstructure:"mongo_url"`
	MongoCollection string `mapstructure:"mongo_collection"`
}
type CacheConfig struct {
	Expiration      time.Duration
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}
