package app

import (
	"log/slog"

	"github.com/google/go-cmp/cmp"
)

type SwitchController struct {
	light Light
}

func NewSwitchController(light Light) *SwitchController {
	return &SwitchController{light: light}
}

func (sc *SwitchController) OnStateChanged(event *StateChangeEvent) error {
	data := event.Event.Data
	v1 := data.OldState
	v2 := data.NewState

	if diff := cmp.Diff(v1, v2); diff != "" {
		slog.Info(diff, "entity", data.EntityID)
	}

	return nil
}
