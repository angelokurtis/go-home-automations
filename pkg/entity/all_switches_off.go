package entity

import (
	"context"
	"github.com/angelokurtis/go-home-automations/pkg/action"
	ga "saml.dev/gome-assistant"
)

type AllSwitchesOff struct {
	act *action.AllSwitchesOff
}

func NewAllSwitchesOff(act *action.AllSwitchesOff) *AllSwitchesOff {
	return &AllSwitchesOff{act: act}
}

func (a *AllSwitchesOff) OnChange(ctx context.Context, entity ga.EntityData) error {
	return a.act.Execute(ctx)
}

func (a *AllSwitchesOff) EntityIds() []string {
	return []string{"input_button.apagar_todas_luzes"}
}
