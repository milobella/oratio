package config

import (
	"encoding/json"
	"time"

	"github.com/milobella/oratio/internal/model"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Server    Server
	Tracing   Tracing
	Auth      Auth
	Cerebro   Cerebro
	Anima     Anima
	Abilities Abilities
}

// fun String() : Serialization function of Config (for logging)
func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		logrus.Fatalf("Configuration serialization error %s", err)
	}
	return string(b)
}

type Server struct {
	ServiceName string `mapstructure:"service_name"`
	Port        int
	LogLevel    string `mapstructure:"log_level"`
}

type Tracing struct {
	ServiceName         string `mapstructure:"service_name"`
	JaegerAgentHostName string `mapstructure:"jaeger_agent_hostname"`
	JaegerAgentPort     int    `mapstructure:"jaeger_agent_port"`
}

type Auth struct {
	AppSecret string `mapstructure:"app_secret"`
}

type Cerebro struct {
	Host               string
	Port               int
	UnderstandEndpoint string `mapstructure:"understand_endpoint"`
}

type Anima struct {
	Host              string
	Port              int
	RestituteEndpoint string `mapstructure:"restitute_endpoint"`
}

type Abilities struct {
	List       []model.Ability
	Cache      Cache
	Database   Database
	StopIntent string `mapstructure:"stop_intent"`
}

type Database struct {
	MongoDatabase   string `mapstructure:"mongo_database"`
	MongoUrl        string `mapstructure:"mongo_url"`
	MongoCollection string `mapstructure:"mongo_collection"`
}
type Cache struct {
	Expiration      time.Duration
	CleanupInterval time.Duration `mapstructure:"cleanup_interval"`
}
