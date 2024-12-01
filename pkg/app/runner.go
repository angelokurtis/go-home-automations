package app

import (
	"context"
	"log/slog"

	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/term"
)

type Runner struct {
	renderer term.Renderer
	ha       HomeAssistant
}

func NewRunner(renderer term.Renderer, ha HomeAssistant) *Runner {
	return &Runner{renderer: renderer, ha: ha}
}

func (r *Runner) Run(ctx context.Context) error {
	eventListener := ga.NewEventListener().
		EventTypes("state_changed").
		Call(func(service *ga.Service, state ga.State, data ga.EventData) {
			err := r.renderer.RenderJSON(data.RawEventJSON)
			if err != nil {
				slog.ErrorContext(ctx, "failed to render json", err)
			}
		}).
		Build()
	r.ha.RegisterEventListeners(eventListener)
	r.ha.Start()

	return nil
}
