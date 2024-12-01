package homeassistant

import (
	env "github.com/caarlos0/env/v11"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type Config struct {
	IpAddress        string
	Port             string `envDefault:"8123"`
	AuthToken        string
	HomeZoneEntityId string
	Secure           bool `envDefault:"false"`
}

func NewConfigFromEnv() (*Config, error) {
	var config Config
	if err := env.ParseWithOptions(&config, env.Options{
		RequiredIfNoDef:       true,
		Prefix:                "HA_",
		UseFieldNameByDefault: true,
	}); err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}
