package action

import (
	"context"

	"github.com/angelokurtis/go-otel/span"
	"github.com/sourcegraph/conc/pool"

	"github.com/angelokurtis/go-home-automations/pkg/app"
)

type EveningLightsOn struct {
	pc app.PowerControl
}

func NewEveningLightsOn(pc app.PowerControl) *EveningLightsOn {
	return &EveningLightsOn{pc: pc}
}

func (e *EveningLightsOn) EntityIDs() []string {
	return []string{
		"light.arandela_da_sala",
		"light.arandela_da_copa",
		"light.numero_da_casa",
		"light.chao_do_jardim",
		"light.lateral",
		"light.lateral2",
	}
}

func (e *EveningLightsOn) LightsOn(ctx context.Context) error {
	p := pool.New().WithMaxGoroutines(10).WithErrors()

	for _, entityID := range e.EntityIDs() {
		entityID := entityID // avoid loop variable capture

		p.Go(func() error {
			ctx, end := span.Start(ctx)
			defer end()

			if err := e.pc.TurnOn(ctx, entityID); err != nil {
				return span.Error(ctx, err)
			}

			return nil
		})
	}

	return p.Wait()
}
