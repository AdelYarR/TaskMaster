package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	BindAddr string `yaml:"bindaddr"`
	Store    struct {
		DBUrl string `mapstructure:"db_url"`
	} `yaml:"store"`
}

func New() (*Config, error) {
	cfg := &Config{}

	viper.SetConfigFile("configs/config/config.yaml")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
