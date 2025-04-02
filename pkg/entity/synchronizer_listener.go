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
	"github.com/angelokurtis/go-home-automations/pkg/app"
)

type SynchronizedSwitchesListeners []*SynchronizedSwitchesListener

func NewSynchronizedSwitchesListeners(pc app.PowerControl, state ga.State) SynchronizedSwitchesListeners {
	listeners := SynchronizedSwitchesListeners{
		// arandelas de dentro
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_6x_da_copa_l1",
			"switch.interruptor_6x_da_copa_l2",
			"switch.interruptor_3x_da_entrada_left",
		),
		// luz da sala
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_3x_da_entrada_center",
			"switch.interruptor_6x_da_entrada_l1",
		),
		// luz da entrada
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_3x_da_entrada_right",
			"switch.interruptor_6x_da_entrada_l2",
		),
		// luz da churrasqueira
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_2x_da_area_gourmet_center",
			"switch.interruptor_6x_da_area_gourmet_l5",
			"switch.interruptor_6x_da_copa_l4",
		),
		// arandelas da piscina
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_6x_da_area_gourmet_l1",
			"switch.interruptor_6x_da_area_gourmet_l6",
		),
		// arandela da lateral
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_2x_da_area_gourmet_left",
			"switch.interruptor_6x_da_copa_l5",
		),
		// luzes do jardim
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_6x_da_entrada_l3",
			"switch.interruptor_6x_da_entrada_l4",
		),
		// luzes do mezanino
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_1x_do_mezanino",
			"switch.interruptor_6x_da_copa_l6",
		),
		// luzes da cozinha
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_da_cozinha_l1",
			"switch.interruptor_da_cozinha2_l1",
		),
		// luzes da lavanderia
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_da_lavanderia_l1",
			"switch.interruptor_da_cozinha_l2",
			"switch.interruptor_da_cozinha2_l2",
		),
		// luzes da copa
		newSynchronizedSwitchesListener(
			pc,
			state,
			"switch.interruptor_6x_da_copa2_l1",
			"switch.interruptor_6x_da_copa2_l2",
		),
	}

	return listeners
}

type SynchronizedSwitchesListener struct {
	pc        app.PowerControl
	state     ga.State
	entityIds []string
	sync.Mutex
}

func newSynchronizedSwitchesListener(pc app.PowerControl, state ga.State, entityIds ...string) *SynchronizedSwitchesListener {
	return &SynchronizedSwitchesListener{pc: pc, state: state, entityIds: entityIds}
}

func (l *SynchronizedSwitchesListener) OnChange(ctx context.Context, entity ga.EntityData) error {
	l.Lock()
	defer l.Unlock()

	state, err := l.state.Get(entity.TriggerEntityId)
	if err != nil {
		return errors.WithStack(err)
	}

	if state.State != entity.ToState {
		slog.DebugContext(ctx, "No action needed",
			slog.String("trigger_entity_id", entity.TriggerEntityId),
			slog.String("trigged_state", entity.ToState),
			slog.String("current_state", state.State),
		)

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
				if err := l.pc.TurnOn(ctx, entityId); err != nil {
					return span.Error(ctx, err)
				}
			}

			if entity.ToState == "off" {
				if err := l.pc.TurnOff(ctx, entityId); err != nil {
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
