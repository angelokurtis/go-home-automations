package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type EventListener interface {
	Call(ctx context.Context, event ga.EventData) error
}
