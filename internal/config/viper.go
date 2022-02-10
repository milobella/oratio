package config

import (
	"github.com/iamolegga/enviper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Read() *Config {
	e := enviper.New(viper.New())
	e.SetEnvPrefix("ORATIO")

	e.SetConfigName("config")
	e.SetConfigType("toml")
	e.AddConfigPath("/etc/oratio/")
	e.AddConfigPath("$HOME/.oratio")
	e.AddConfigPath(".")
	err := e.ReadInConfig()
	if err != nil {
		fatal(err)
	}

	var config Config
	if err = e.Unmarshal(&config); err != nil {
		fatal(err)
	} else {
		logrus.Info("Successfully red configuration !")
		logrus.Debugf("-> %+v", config)
	}

	logrus.SetLevel(config.Server.LogLevel)

	return &config
}

func fatal(err error) {
	logrus.WithError(err).Fatal("Error reading config.")
}
