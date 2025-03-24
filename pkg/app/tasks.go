package app

import (
	"context"
)

type DailyTask interface {
	Execute(ctx context.Context) error
	ScheduledTime() string
	Name() string
}

type SunsetTask interface {
	Execute(ctx context.Context) error
	SunsetOffset() string
	Name() string
}

type SunriseTask interface {
	Execute(ctx context.Context) error
	SunriseOffset() string
	Name() string
}
