package configs

import (
	"fmt"
	"github.com/vrischmann/envconfig"
)

type Config struct {
	Port int `envconfig:"PORT,default=8080"`
}

func NewConfig() (*Config, error) {
	config := &Config{}
	if err := envconfig.Init(config); err != nil {
		return nil, fmt.Errorf("Field cread new config: %w\n", err)
	}
	return config, nil
}
