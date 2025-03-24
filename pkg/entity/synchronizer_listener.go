package entity

import (
	"context"
	"log/slog"
	"sync"

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
	sync.Mutex
}

func NewSynchronizedSwitchesListener(service *ga.Service, state ga.State, entityIds ...string) *SynchronizedSwitchesListener {
	return &SynchronizedSwitchesListener{service: service, state: state, entityIds: entityIds}
}

func (l *SynchronizedSwitchesListener) OnChange(ctx context.Context, entity ga.EntityData) error {
	l.Lock()
	defer l.Unlock()

	state, err := l.state.Get(entity.TriggerEntityId)
	if err != nil {
		return errors.WithStack(err)
	}

	if state.State != entity.ToState {
		if state.State != entity.ToState {
			slog.DebugContext(ctx, "State is already set, no action needed",
				slog.String("trigger_entity_id", entity.TriggerEntityId),
				slog.String("current_state", state.State),
				slog.String("new_state", entity.ToState),
			)

			return nil
		}

		return nil
	}

	entityIds := lo.Without(l.entityIds, entity.TriggerEntityId)
	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for _, entityId := range entityIds {
		entityId := entityId // avoid loop variable capture

		p.Go(func() error {
			ctx, end := span.Start(ctx)
			defer end()

			if entity.ToState == "on" {
				if err := l.turnOn(ctx, entityId); err != nil {
					return span.Error(ctx, err)
				}
			}

			if entity.ToState == "off" {
				if err := l.turnOff(ctx, entityId); err != nil {
					return span.Error(ctx, err)
				}
			}

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
		return errors.WithStack(err)
	}

	if state.State == "on" {
		slog.DebugContext(ctx, "Entity already on",
			slog.String("entity_id", entityId),
		)

		return nil
	}

	if err = l.service.HomeAssistant.TurnOn(entityId); err != nil {
		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned on",
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
		slog.DebugContext(ctx, "Entity already off",
			slog.String("entity_id", entityId),
		)

		return nil
	}

	if err = l.service.HomeAssistant.TurnOff(entityId); err != nil {
		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned off",
		slog.String("entity_id", entityId),
	)

	return nil
}
