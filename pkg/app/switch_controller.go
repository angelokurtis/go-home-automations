package app

import (
	"log/slog"

	"github.com/google/go-cmp/cmp"
)

type SwitchController struct {
	light Light
	state State
}

func NewSwitchController(light Light, state State) *SwitchController {
	return &SwitchController{light: light, state: state}
}

func (sc *SwitchController) OnStateChanged(event *StateChangeEvent) error {
	entityID := event.EntityID

	diff := cmp.Diff(event.OldState, event.NewState)
	if diff == "" {
		return nil
	}

	slog.Info(diff, "entity", entityID)

	return nil
}
