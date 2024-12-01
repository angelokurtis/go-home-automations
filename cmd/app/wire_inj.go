package main

import (
	"context"

	"github.com/google/wire"
)

type AppRunner interface {
	Run(ctx context.Context) error
}

func newAppRunner(ctx context.Context) (AppRunner, func(), error) {
	wire.Build(Providers)
	return nil, nil, nil
}
