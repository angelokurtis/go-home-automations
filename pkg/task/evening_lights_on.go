package task

import (
	"context"
	"log/slog"

	"github.com/angelokurtis/go-otel/span"
	"github.com/sourcegraph/conc/pool"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type EveningLightsOn struct {
	service *ga.Service
	state   ga.State
}

func NewEveningLightsOn(service *ga.Service, state ga.State) *EveningLightsOn {
	return &EveningLightsOn{service: service, state: state}
}

func (e *EveningLightsOn) EntityIDs() []string {
	return []string{
		"light.arandela_da_sala",
		"light.arandela_da_copa",
		"light.numero_da_casa",
		"light.chao_do_jardim",
		"light.lateral",
	}
}

func (e *EveningLightsOn) Execute(ctx context.Context) error {
	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for _, entityID := range e.EntityIDs() {
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

func (e *EveningLightsOn) SunsetOffset() string {
	return "1m"
}

func (e *EveningLightsOn) Name() string {
	return "Luzes ao Pôr do Sol"
}

func (e *EveningLightsOn) turnOn(ctx context.Context, entityId string, serviceData ...map[string]any) error {
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
