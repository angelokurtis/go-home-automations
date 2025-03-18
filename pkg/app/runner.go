package app

import (
	"context"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
)

type Runner struct {
	app *ga.App
}

func NewRunner(app *ga.App) *Runner {
	return &Runner{app: app}
}

func (r *Runner) Run(ctx context.Context) error {
	return errors.New("not implemented")
}
