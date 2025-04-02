package action

import (
	"context"
	"strings"

	"github.com/angelokurtis/go-otel/span"
	"github.com/samber/lo"
	"github.com/sourcegraph/conc/pool"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
	"github.com/angelokurtis/go-home-automations/pkg/app"
)

type AllSwitchesOff struct {
	er app.EntityReader
	pc app.PowerControl
}

func NewAllSwitchesOff(er app.EntityReader, pc app.PowerControl) *AllSwitchesOff {
	return &AllSwitchesOff{er: er, pc: pc}
}

func (a *AllSwitchesOff) Execute(ctx context.Context) error {
	entities, err := a.er.ListEntities()
	if err != nil {
		return errors.WithStack(err)
	}

	entities = lo.Filter(entities, func(item ga.EntityState, index int) bool {
		return strings.HasPrefix(item.EntityID, "switch.") && item.EntityID != "switch.zigbee2mqtt_bridge_permit_join"
	})

	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for _, entity := range entities {
		entity := entity // avoid loop variable capture

		p.Go(func() error {
			ctx, end := span.Start(ctx)
			defer end()

			if err := a.pc.TurnOff(ctx, entity.EntityID); err != nil {
				return span.Error(ctx, err)
			}

			return nil
		})
	}

	return p.Wait()
}
