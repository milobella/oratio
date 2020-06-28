package config

type ServerConfiguration struct {
	Port     int    `env:"SERVER_PORT" default:"8080"`
	LogLevel string `id:"log_level" env:"SERVER_LOG_LEVEL" default:"<root>=ERROR"`
}

type CerebroConfiguration struct {
	Host string `env:"CEREBRO_HOST" default:"localhost"`
	Port int    `env:"CEREBRO_PORT" default:"9444"`
}

type AnimaConfiguration struct {
	Host string `env:"ANIMA_HOST" default:"localhost"`
	Port int    `env:"ANIMA_PORT" default:"9333"`
}

type AbilityMongoConfiguration struct {
	Url string `env:"ABILITY_MONGO_URL"`
}

type AbilityConfiguration struct {
	Name    string
	Host    string
	Port    int
	Intents []string
}
