package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type DailyTask interface {
	Execute(ctx context.Context, state ga.State) error
	ScheduledTime() string
}

type SunsetTask interface {
	Execute(ctx context.Context, state ga.State) error
	SunsetOffset() string
}

type SunriseTask interface {
	Execute(ctx context.Context, state ga.State) error
	SunriseOffset() string
}
