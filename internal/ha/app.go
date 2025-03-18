package ha

import (
	"context"
	"log/slog"

	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

func NewApp(ctx context.Context, config *Config) (*ga.App, func(), error) {
	app, err := ga.NewApp(ga.NewAppRequest{
		URL:              config.URL,
		HAAuthToken:      config.AuthToken,
		HomeZoneEntityId: config.HomeZoneEntityId,
	})
	if err != nil {
		return nil, func() {}, errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Home Assistant App initialized",
		slog.String("url", config.URL),
		slog.String("home_zone_entity_id", config.HomeZoneEntityId),
	)

	return app, func() {
		app.Cleanup()
		slog.InfoContext(ctx, "Home Assistant App cleanup completed",
			slog.String("url", config.URL),
			slog.String("home_zone_entity_id", config.HomeZoneEntityId),
		)
	}, nil
}
