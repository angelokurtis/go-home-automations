package app

import (
	"context"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type Runner struct{}

func (r *Runner) Run(ctx context.Context) error {
	return errors.New("not implemented")
}
