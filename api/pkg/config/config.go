package config

import (
	"github.com/caarlos0/env/v11"
)

func ReadEnvConfig(cfg any) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}
	return nil
}
