package main

import (
	"github.com/google/wire"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/ha"
	"github.com/angelokurtis/go-home-automations/internal/term"
	"github.com/angelokurtis/go-home-automations/pkg/app"
	"github.com/angelokurtis/go-home-automations/pkg/entity"
)

//nolint:unused // This function is used during compile-time to generate code for dependency injection
var providers = wire.NewSet(
	app.NewRunner,
	ha.Providers,
	term.Providers,
	wire.Bind(new(Runner), new(*app.Runner)),
	wire.Bind(new(term.Renderer), new(*term.MarkdownRenderer)),

	entityListeners,
	eventListeners,
)

func entityListeners(service *ga.Service) []app.EntityListener {
	return []app.EntityListener{
		// arandelas
		entity.NewSynchronizedSwitchesListener(
			service,
			"switch.interruptor_6x_da_copa_l1",
			"switch.interruptor_6x_da_copa_l2",
		),
	}
}

func eventListeners() []app.EventListener {
	return []app.EventListener{}
}
