package ha

import (
	envv11 "github.com/caarlos0/env/v11"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type Config struct {
	Home string `env:"HOME"`
}

func LoadConfig() (*Config, error) {
	cfg, err := envv11.ParseAs[Config]()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &cfg, nil
}
