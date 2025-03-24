package app

import (
	"context"
)

type DailyTask interface {
	Execute(ctx context.Context) error
	ScheduledTime() string
}

type SunsetTask interface {
	Execute(ctx context.Context) error
	SunsetOffset() string
}

type SunriseTask interface {
	Execute(ctx context.Context) error
	SunriseOffset() string
}
