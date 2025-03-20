package app

import (
	"context"
	"log/slog"

	"github.com/angelokurtis/go-otel/span"
	"github.com/lmittmann/tint"
	"github.com/samber/lo"
	ga "saml.dev/gome-assistant"
)

const (
	arandelaDaSala = "switch.interruptor_6x_da_copa_l1"
	arandelaDaCopa = "switch.interruptor_6x_da_copa_l2"
)

type Runner struct {
	app             *ga.App
	entityListeners []EntityListener
	eventListeners  []EventListener
}

func NewRunner(app *ga.App, entityListeners []EntityListener, eventListeners []EventListener) *Runner {
	return &Runner{app: app, entityListeners: entityListeners, eventListeners: eventListeners}
}

func (r *Runner) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting application")

	entityListeners := lo.Map(r.entityListeners, func(entityListener EntityListener, index int) ga.EntityListener {
		return ga.NewEntityListener().
			EntityIds(entityListener.EntityIds()...).
			Call(func(service *ga.Service, state ga.State, entity ga.EntityData) {
				slog.InfoContext(ctx, "Entity state changed",
					slog.String("entity_id", entity.TriggerEntityId),
					slog.String("new_state", entity.ToState),
				)

				if err := entityListener.OnChange(ctx, entity); err != nil {
					slog.ErrorContext(ctx, "Failed to process entity change", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "Entity change processed")
			}).
			Build()
	})
	r.app.RegisterEntityListeners(entityListeners...)
	slog.InfoContext(context.Background(), "Registered entity listeners",
		slog.Int("count", len(entityListeners)),
	)

	eventListeners := lo.Map(r.eventListeners, func(eventListener EventListener, index int) ga.EventListener {
		return ga.NewEventListener().
			EventTypes(eventListener.EventTypes()...).
			Call(func(service *ga.Service, state ga.State, event ga.EventData) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "Event received",
					slog.String("event_type", event.Type),
					slog.String("raw", string(event.RawEventJSON)),
				)

				if err := eventListener.OnEvent(ctx, event); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "Failed to process event", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "Event processed")
			}).
			Build()
	})
	r.app.RegisterEventListeners(eventListeners...)
	slog.InfoContext(context.Background(), "Registered event listeners",
		slog.Int("count", len(eventListeners)),
	)

	r.app.Start()

	slog.InfoContext(ctx, "Application started")

	return nil
}
