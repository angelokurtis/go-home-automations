package app

import (
	"log/slog"
	"strings"

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
	data := event.Event.Data

	entityID := data.EntityID
	if !strings.HasPrefix(entityID, "switch.") {
		return nil
	}

	v1 := data.OldState
	v2 := data.NewState

	diff := cmp.Diff(v1.State, v2.State)
	if diff == "" {
		return nil
	}

	slog.Info(diff, "entity", entityID)

	return nil
}
