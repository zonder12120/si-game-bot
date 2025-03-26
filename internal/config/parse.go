package config

import (
	"github.com/caarlos0/env/v11"
)

func Parse() (App, error) {
	var cfg App

	err := env.Parse(&cfg)
	if err != nil {
		return App{}, err
	}

	return cfg, nil
}
