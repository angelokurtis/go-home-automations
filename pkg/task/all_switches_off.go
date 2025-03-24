package task

import (
	"context"
	"log/slog"
	"strings"

	"github.com/angelokurtis/go-otel/span"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type AllSwitchesOff struct {
	service *ga.Service
	state   ga.State
}

func NewAllSwitchesOff(service *ga.Service, state ga.State) *AllSwitchesOff {
	return &AllSwitchesOff{service: service, state: state}
}

func (a *AllSwitchesOff) Execute(ctx context.Context) error {
	entities, err := a.state.ListEntities()
	if err != nil {
		return errors.WithStack(err)
	}

	entities = lo.Filter(entities, func(item ga.EntityState, index int) bool {
		return strings.HasPrefix(item.EntityID, "switch.") && item.EntityID != "switch.zigbee2mqtt_bridge_permit_join"
	})

	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for _, entity := range entities {
		entity := entity

		p.Go(func() error {
			ctx, end := span.Start(ctx)
			defer end()

			if err := a.turnOff(ctx, entity.EntityID); err != nil {
				return span.Error(ctx, err)
			}

			return nil
		})
	}

	return p.Wait()
}

func (a *AllSwitchesOff) ScheduledTime() string {
	return "17:04"
}

func (a *AllSwitchesOff) Name() string {
	return "Apagar Todas os Interruptores"
}

func (a *AllSwitchesOff) turnOff(ctx context.Context, entityId string, serviceData ...map[string]any) error {
	state, err := a.state.Get(entityId)
	if err != nil {
		return errors.WithStack(err)
	}

	if state.State == "off" {
		slog.DebugContext(ctx, "Entity already off",
			slog.String("entity_id", entityId),
		)

		return nil
	}

	if err = a.service.HomeAssistant.TurnOff(entityId); err != nil {
		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned off",
		slog.String("entity_id", entityId),
	)

	return nil
}
