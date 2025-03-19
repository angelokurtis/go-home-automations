package ha

import (
	envv11 "github.com/caarlos0/env/v11"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type Config struct {
	URL              string `env:"URL" envDefault:"http://192.168.1.123:8123"`
	AuthToken        string `env:"AUTH_TOKEN"`
	HomeZoneEntityId string `env:"HOME_ZONE_ENTITY_ID" envDefault:"zone.home"`
}

func LoadConfig() (*Config, error) {
	config, err := envv11.ParseAsWithOptions[Config](envv11.Options{
		Prefix: "HA_",
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &config, nil
}
