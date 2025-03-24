package task

import (
	"context"
	"log/slog"

	"github.com/angelokurtis/go-otel/span"
	"github.com/sourcegraph/conc/pool"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type EveningLights struct {
	service   *ga.Service
	state     ga.State
	entityIDs []string
}

func NewEveningLights(service *ga.Service, state ga.State, entityIDs []string) *EveningLights {
	return &EveningLights{service: service, state: state, entityIDs: entityIDs}
}

func (e *EveningLights) Execute(ctx context.Context) error {
	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for _, entityID := range e.entityIDs {
		entityID := entityID // avoid loop variable capture

		p.Go(func() error {
			ctx, end := span.Start(ctx)
			defer end()

			if err := e.turnOn(ctx, entityID); err != nil {
				return span.Error(ctx, err)
			}

			return nil
		})
	}

	return p.Wait()
}

func (e *EveningLights) SunsetOffset() string {
	return "1m"
}

func (e *EveningLights) Name() string {
	return "Luzes ao Pôr do Sol"
}

func (e *EveningLights) turnOn(ctx context.Context, entityId string, serviceData ...map[string]any) error {
	state, err := e.state.Get(entityId)
	if err != nil {
		return errors.WithStack(err)
	}

	if state.State == "on" {
		slog.DebugContext(ctx, "Entity already on",
			slog.String("entity_id", entityId),
		)

		return nil
	}

	if err = e.service.HomeAssistant.TurnOn(entityId); err != nil {
		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned on",
		slog.String("entity_id", entityId),
	)

	return nil
}
