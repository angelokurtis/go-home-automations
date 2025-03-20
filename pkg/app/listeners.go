package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type EntityListener interface {
	OnChange(ctx context.Context, entity ga.EntityData) error
	EntityIds() []string
}

type EventListener interface {
	OnEvent(ctx context.Context, event ga.EventData) error
	EventTypes() []string
}
