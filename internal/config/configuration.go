package config

type ServerConfiguration struct {
	Port int 		`env:"SERVER_PORT" default:"8080"`
}

type CerebroConfiguration struct {
	Host string 	`env:"CEREBRO_HOST" default:"localhost"`
	Port int 		`env:"CEREBRO_PORT" default:"9444"`
}

type AnimaConfiguration struct {
	Host string 	`env:"ANIMA_HOST" default:"localhost"`
	Port int 		`env:"ANIMA_PORT" default:"9333"`
}

type AbilityConfiguration struct {
	Intents []string
	Host string
	Port int
}
