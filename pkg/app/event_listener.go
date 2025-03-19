package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type EventListener interface {
	OnEvent(ctx context.Context, event ga.EventData) error
}
