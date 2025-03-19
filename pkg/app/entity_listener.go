package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type EntityListener interface {
	OnChange(ctx context.Context, entity ga.EntityData) error
	EntityIds() []string
}
