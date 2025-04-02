package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type PowerControl interface {
	TurnOn(ctx context.Context, entityID string, serviceData ...map[string]any) error
	TurnOff(ctx context.Context, entityID string) error
}

type EntityReader interface {
	ListEntities() ([]ga.EntityState, error)
	Get(entityId string) (ga.EntityState, error)
}
