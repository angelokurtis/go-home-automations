package entity

import (
	"context"
	"log/slog"

	"github.com/angelokurtis/go-otel/span"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type SynchronizedSwitchesListener struct {
	service   *ga.Service
	state     ga.State
	entityIds []string
}

func NewSynchronizedSwitchesListener(service *ga.Service, state ga.State, entityIds ...string) *SynchronizedSwitchesListener {
	return &SynchronizedSwitchesListener{service: service, state: state, entityIds: entityIds}
}

func (l *SynchronizedSwitchesListener) OnChange(ctx context.Context, entity ga.EntityData) error {
	entityIds := lo.Without(l.entityIds, entity.TriggerEntityId)
	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for _, entityId := range entityIds {
		entityId := entityId // avoid loop variable capture

		p.Go(func() error {
			ctx, end := span.Start(ctx)
			defer end()

			if entity.ToState == "on" {
				if err := l.service.HomeAssistant.TurnOn(entityId); err != nil {
					return errors.WithStack(span.Error(ctx, err))
				}
			}

			if entity.ToState == "off" {
				if err := l.service.HomeAssistant.TurnOff(entityId); err != nil {
					return errors.WithStack(span.Error(ctx, err))
				}
			}

			slog.InfoContext(ctx, "Entity state updated",
				slog.String("entity-id", entityId),
				slog.String("new-state", entity.ToState),
			)

			return nil
		})
	}

	return p.Wait()
}

func (l *SynchronizedSwitchesListener) EntityIds() []string {
	return l.entityIds
}

func (l *SynchronizedSwitchesListener) turnOn(ctx context.Context, entityId string, serviceData ...map[string]any) error {
	state, err := l.state.Get(entityId)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to get entity state",
			slog.String("entity_id", entityId),
			slog.Any("error", err),
		)

		return errors.WithStack(err)
	}

	if state.State == "on" {
		slog.InfoContext(ctx, "Entity already turned on",
			slog.String("entity_id", entityId),
		)

		return nil
	}

	if err = l.service.HomeAssistant.TurnOn(entityId); err != nil {
		slog.ErrorContext(ctx, "Failed to turn on entity",
			slog.String("entity_id", entityId),
			slog.Any("error", err),
		)

		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned on successfully",
		slog.String("entity_id", entityId),
	)

	return nil
}

func (l *SynchronizedSwitchesListener) turnOff(ctx context.Context, entityId string, serviceData ...map[string]any) error {
	state, err := l.state.Get(entityId)
	if err != nil {
		return errors.WithStack(err)
	}

	if state.State == "off" {
		slog.InfoContext(ctx, "Entity already turned off",
			slog.String("entity_id", entityId),
		)

		return nil
	}

	if err = l.service.HomeAssistant.TurnOff(entityId); err != nil {
		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned off successfully",
		slog.String("entity_id", entityId),
	)

	return nil
}
