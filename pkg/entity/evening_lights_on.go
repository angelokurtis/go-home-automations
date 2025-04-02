package entity

import (
	"context"
	"github.com/angelokurtis/go-home-automations/pkg/action"
	ga "saml.dev/gome-assistant"
)

type EveningLightsOn struct {
	act action.EveningLightsOn
}

func (e *EveningLightsOn) OnChange(ctx context.Context, entity ga.EntityData) error {
	return e.act.LightsOn(ctx)
}

func (e *EveningLightsOn) EntityIds() []string {
	return []string{"input_button.luzes_do_entardecer"}
}
