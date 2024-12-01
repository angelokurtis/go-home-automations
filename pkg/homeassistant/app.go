package homeassistant

import (
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

func NewApp(config *Config) (*ga.App, func(), error) {
	app, err := ga.NewApp(ga.NewAppRequest{
		IpAddress:        config.IpAddress,
		Port:             config.Port,
		HAAuthToken:      config.AuthToken,
		HomeZoneEntityId: config.HomeZoneEntityId,
		Secure:           config.Secure,
	})
	if err != nil {
		return nil, func() {}, errors.WithStack(err)
	}

	return app, app.Cleanup, nil
}
