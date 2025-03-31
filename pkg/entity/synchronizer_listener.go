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

type SynchronizedSwitchesListeners []*SynchronizedSwitchesListener

func NewSynchronizedSwitchesListeners(service *ga.Service, state ga.State) SynchronizedSwitchesListeners {
	listeners := SynchronizedSwitchesListeners{
		// arandelas de dentro
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_6x_da_copa_l1",
			"switch.interruptor_6x_da_copa_l2",
			"switch.interruptor_3x_da_entrada_left",
		),
		// luz da sala
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_3x_da_entrada_center",
			"switch.interruptor_6x_da_entrada_l1",
		),
		// luz da entrada
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_3x_da_entrada_right",
			"switch.interruptor_6x_da_entrada_l2",
		),
		// luz da churrasqueira
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_2x_da_area_gourmet_center",
			"switch.interruptor_6x_da_area_gourmet_l5",
			"switch.interruptor_6x_da_copa_l4",
		),
		// arandelas da piscina
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_6x_da_area_gourmet_l1",
			"switch.interruptor_6x_da_area_gourmet_l6",
		),
		// arandela da lateral
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_2x_da_area_gourmet_left",
			"switch.interruptor_6x_da_copa_l5",
		),
		// luzes do jardim
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_6x_da_entrada_l3",
			"switch.interruptor_6x_da_entrada_l4",
		),
		// luzes do mezanino
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_1x_do_mezanino",
			"switch.interruptor_6x_da_copa_l6",
		),
		// luzes da cozinha
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_da_cozinha_l1",
			"switch.interruptor_da_cozinha2_l1",
		),
		// luzes da lavanderia
		newSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_da_lavanderia_l1",
			"switch.interruptor_da_cozinha_l2",
			"switch.interruptor_da_cozinha2_l2",
		),
	}

	return listeners
}

type SynchronizedSwitchesListener struct {
	service   *ga.Service
	state     ga.State
	entityIds []string
	sync.Mutex
}

func newSynchronizedSwitchesListener(service *ga.Service, state ga.State, entityIds ...string) *SynchronizedSwitchesListener {
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
