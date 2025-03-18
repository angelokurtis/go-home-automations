package ha

import (
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

func NewApp(config *Config) (*ga.App, error) {
	app, err := ga.NewApp(ga.NewAppRequest{
		URL:              config.URL,
		HAAuthToken:      config.AuthToken,
		HomeZoneEntityId: config.HomeZoneEntityId,
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return app, nil
}
