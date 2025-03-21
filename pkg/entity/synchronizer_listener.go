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
	entityIds []string
}

func NewSynchronizedSwitchesListener(service *ga.Service, entityIds ...string) *SynchronizedSwitchesListener {
	return &SynchronizedSwitchesListener{service: service, entityIds: entityIds}
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

func (l *SynchronizedSwitchesListener) turnOn(entityId string, serviceData ...map[string]any) []string {
	// TODO: implement me
	panic("implement me")
}

func (l *SynchronizedSwitchesListener) turnOff(entityId string, serviceData ...map[string]any) []string {
	// TODO: implement me
	panic("implement me")
}
