package app

import (
	"context"
	"log/slog"
	"strings"

	"github.com/buger/jsonparser"
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
			event := new(StateChangeEvent)

			// Define the keys to extract from the JSON event data.
			paths := [][]string{
				{"event", "event_type"},
				{"event", "data", "entity_id"},
				{"event", "data", "old_state", "state"},
				{"event", "data", "new_state", "state"},
			}

			// Extract values from JSON data and map them to the StateChangeEvent struct.
			jsonparser.EachKey(data.RawEventJSON, func(idx int, value []byte, vt jsonparser.ValueType, err error) {
				if err != nil {
					return
				}
				switch idx {
				case 0:
					event.Type = string(value)
				case 1:
					event.EntityID = string(value)
				case 2:
					event.OldState = string(value)
				case 3:
					event.NewState = string(value)
				}
			}, paths...)

			// Filter events: only continue if it's a state change for a switch or light entity.
			if event.Type != "state_changed" ||
				!(strings.HasPrefix(event.EntityID, "switch.") || strings.HasPrefix(event.EntityID, "light.")) {
				return
			}

			// Call the switch controller's OnStateChanged method to handle the event.
			if err := r.switchController.OnStateChanged(event); err != nil {
				slog.Error("error during switch controller execution", tint.Err(err))
			}
		}).
		Build()
	r.ha.RegisterEventListeners(eventListener)
	r.ha.Start()

	return nil
}
