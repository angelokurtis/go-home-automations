//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
)

type Runner interface {
	Run(ctx context.Context) error
}

func NewRunner(ctx context.Context) (Runner, func(), error) {
	wire.Build(providers)
	return nil, nil, nil
}
