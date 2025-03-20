//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/angelokurtis/go-home-automations/internal/ha"
	"github.com/angelokurtis/go-home-automations/pkg/app"
)

import (
	_ "github.com/angelokurtis/go-home-automations/internal/logger"
)

type Runner interface {
	Run(ctx context.Context) error
}

func NewRunner(ctx context.Context) (Runner, func(), error) {
	wire.Build(providers)
	return nil, nil, nil
}
