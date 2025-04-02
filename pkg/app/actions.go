package app

import "context"

type PowerControl interface {
	TurnOn(ctx context.Context, entityID string, data map[string]any) error
	TurnOff(ctx context.Context, entityID string) error
}
