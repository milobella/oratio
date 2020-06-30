package config

import (
	"github.com/iamolegga/enviper"
	"github.com/spf13/viper"
)

func ReadConfiguration() (*Configuration, error) {
	e := enviper.New(viper.New())
	e.SetEnvPrefix("ORATIO")

	e.SetConfigName("config")
	e.SetConfigType("toml")
	e.AddConfigPath("/etc/oratio/")
	e.AddConfigPath("$HOME/.oratio")
	e.AddConfigPath(".")
	err := e.ReadInConfig()
	if err != nil {
		return nil, err
	}

	var C Configuration
	err = e.Unmarshal(&C)
	return &C, err
}
