package main

import (
	"github.com/google/wire"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/pkg/app"
	"github.com/angelokurtis/go-home-automations/pkg/homeassistant"
)

var Providers = wire.NewSet(
	wire.Bind(new(app.HomeAssistant), new(*ga.App)),
	wire.Bind(new(AppRunner), new(*app.Runner)),
	app.Providers,
	homeassistant.Providers,
)
