package app

import (
	"log/slog"
	"strings"

	"github.com/google/go-cmp/cmp"
)

type SwitchController struct {
	light Light
	State State
}

func NewSwitchController(light Light, state State) *SwitchController {
	return &SwitchController{light: light, State: state}
}

func (sc *SwitchController) OnStateChanged(event *StateChangeEvent) error {
	data := event.Event.Data

	entityID := data.EntityID
	if !strings.HasPrefix(entityID, "switch.") {
		return nil
	}

	v1 := data.OldState
	v2 := data.NewState

	if diff := cmp.Diff(v1, v2); diff != "" {
		slog.Info(diff, "entity", entityID)
	}

	return nil
}
