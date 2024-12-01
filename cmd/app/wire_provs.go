package main

import (
	"github.com/google/wire"

	"github.com/angelokurtis/go-home-automations/pkg/app"
)

var Providers = wire.NewSet(
	wire.Bind(new(AppRunner), new(*app.Runner)),
	app.Providers,
)
