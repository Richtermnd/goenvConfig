package config

import "github.com/Richtermnd/goenvConfig/internal/config"

func LoadEnviroment(envFiles ...string) error {
	return config.LoadEnviroment(envFiles...)
}

func LoadConfig(cfg interface{}) error {
	return config.LoadConfig(cfg)
}
