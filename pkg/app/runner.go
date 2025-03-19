package app

import (
	"context"
	"log/slog"

	ga "saml.dev/gome-assistant"
)

const (
	arandelaDaSala = "switch.interruptor_6x_da_copa_l1"
	arandelaDaCopa = "switch.interruptor_6x_da_copa_l2"
)

type Runner struct {
	app *ga.App
}

func NewRunner(app *ga.App) *Runner {
	return &Runner{app: app}
}

func (r *Runner) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting application")

	entityListener := ga.NewEntityListener().
		EntityIds(arandelaDaSala, arandelaDaCopa).
		Call(func(service *ga.Service, state ga.State, sensor ga.EntityData) {
			slog.InfoContext(ctx, "Entity state changed",
				slog.String("entity_id", sensor.TriggerEntityId),
				slog.String("new_state", sensor.ToState),
			)
			if sensor.ToState == "on" {
				err := service.HomeAssistant.TurnOn(arandelaDaSala)
				err := service.HomeAssistant.TurnOn(arandelaDaCopa)
			} else {
				err := service.HomeAssistant.TurnOff(arandelaDaSala)
				err := service.HomeAssistant.TurnOff(arandelaDaCopa)
			}
		}).
		Build()

	r.app.RegisterEntityListeners(entityListener)
	slog.InfoContext(ctx, "Entity listeners registered")

	r.app.Start()

	slog.InfoContext(ctx, "Application started")

	return nil
}
