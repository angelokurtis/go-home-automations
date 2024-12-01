package app

import (
	"context"
)

type Runner struct {
	ha HomeAssistant
}

func NewRunner(ha HomeAssistant) *Runner {
	return &Runner{ha: ha}
}

func (r *Runner) Run(ctx context.Context) error {
	// TODO implement me
	panic("implement me")
}
