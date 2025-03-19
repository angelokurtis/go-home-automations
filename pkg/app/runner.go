package app

import (
	"context"
	"log/slog"

	"github.com/lmittmann/tint"
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
				slog.InfoContext(ctx, "Turning on lights",
					slog.Any("entity_ids", []string{arandelaDaSala, arandelaDaCopa}),
				)
				if err := service.HomeAssistant.TurnOn(arandelaDaSala); err != nil {
					slog.ErrorContext(ctx, "Failed to turn on light",
						slog.String("entity_id", arandelaDaSala),
						tint.Err(err),
					)
					return
				}
				if err := service.HomeAssistant.TurnOn(arandelaDaCopa); err != nil {
					slog.ErrorContext(ctx, "Failed to turn on light",
						slog.String("entity_id", arandelaDaCopa),
						tint.Err(err),
					)
					return
				}
			} else {
				slog.InfoContext(ctx, "Turning off lights",
					slog.Any("entity_ids", []string{arandelaDaSala, arandelaDaCopa}),
				)
				if err := service.HomeAssistant.TurnOff(arandelaDaSala); err != nil {
					slog.ErrorContext(ctx, "Failed to turn off light",
						slog.String("entity_id", arandelaDaSala),
						tint.Err(err),
					)
					return
				}
				if err := service.HomeAssistant.TurnOff(arandelaDaCopa); err != nil {
					slog.ErrorContext(ctx, "Failed to turn off light",
						slog.String("entity_id", arandelaDaCopa),
						tint.Err(err),
					)
					return
				}
			}
		}).
		Build()

	ga.NewEventListener().EventTypes().Call(func(service *ga.Service, state ga.State, data ga.EventData) {
	})

	r.app.RegisterEntityListeners(entityListener)
	slog.InfoContext(ctx, "Entity listeners registered")

	r.app.Start()

	slog.InfoContext(ctx, "Application started")

	return nil
}
