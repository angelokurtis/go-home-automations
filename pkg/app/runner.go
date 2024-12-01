package app

import (
	"context"
)

type Runner struct{}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Run(ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}
