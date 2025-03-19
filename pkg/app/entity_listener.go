package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type EntityListener interface {
	Call(ctx context.Context, entity ga.EntityData) error
}
