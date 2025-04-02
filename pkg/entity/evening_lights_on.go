package entity

import (
	"context"

	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/pkg/action"
)

type EveningLightsOn struct {
	act *action.EveningLightsOn
}

func NewEveningLightsOn(act *action.EveningLightsOn) *EveningLightsOn {
	return &EveningLightsOn{act: act}
}

func (e *EveningLightsOn) OnChange(ctx context.Context, entity ga.EntityData) error {
	return e.act.LightsOn(ctx)
}

func (e *EveningLightsOn) EntityIds() []string {
	return []string{"input_button.luzes_do_entardecer"}
}
