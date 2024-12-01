package app

import (
	"context"
	"log/slog"

	"github.com/lmittmann/tint"
	ga "saml.dev/gome-assistant"
)

type Runner struct {
	ha               HomeAssistant
	switchController *SwitchController
}

func NewRunner(ha HomeAssistant, switchController *SwitchController) *Runner {
	return &Runner{ha: ha, switchController: switchController}
}

func (r *Runner) Run(ctx context.Context) error {
	eventListener := ga.NewEventListener().
		EventTypes("state_changed").
		Call(func(service *ga.Service, state ga.State, data ga.EventData) {
			event, err := UnmarshalStateChangeEvent(data.RawEventJSON)
			if err != nil {
				slog.Warn("failed to unmarshal state change event")
				return
			}

			if err = r.switchController.OnStateChanged(event); err != nil {
				slog.Error("error during switch controller execution", tint.Err(err))
			}
		}).
		Build()
	r.ha.RegisterEventListeners(eventListener)
	r.ha.Start()

	return nil
}
