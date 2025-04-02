package action

import (
	"context"
	"log/slog"

	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type PowerControl struct {
	svc   *ga.Service
	state ga.State
}

func NewPowerControl(svc *ga.Service, state ga.State) *PowerControl {
	return &PowerControl{svc: svc, state: state}
}

func (p *PowerControl) TurnOn(ctx context.Context, entityID string, serviceData ...map[string]any) error {
	state, err := p.state.Get(entityID)
	if err != nil {
		return errors.WithStack(err)
	}

	if state.State == "on" {
		slog.DebugContext(ctx, "Entity already on", slog.String("entity_id", entityID))
		return nil
	}

	if err = p.svc.HomeAssistant.TurnOn(entityID, serviceData...); err != nil {
		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned on", slog.String("entity_id", entityID))

	return nil
}

func (p *PowerControl) TurnOff(ctx context.Context, entityID string) error {
	state, err := p.state.Get(entityID)
	if err != nil {
		return errors.WithStack(err)
	}

	if state.State == "off" {
		slog.DebugContext(ctx, "Entity already off", slog.String("entity_id", entityID))
		return nil
	}

	if err = p.svc.HomeAssistant.TurnOff(entityID); err != nil {
		return errors.WithStack(err)
	}

	slog.InfoContext(ctx, "Entity turned off", slog.String("entity_id", entityID))

	return nil
}
