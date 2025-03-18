package app

import (
	"context"

	ga "saml.dev/gome-assistant"
)

type Runner struct {
	app *ga.App
}

func NewRunner(app *ga.App) *Runner {
	return &Runner{app: app}
}

func (r *Runner) Run(ctx context.Context) error {
	entityListener := ga.NewEntityListener().
		EntityIds("binary_sensor.pantry_door").
		Call(func(service *ga.Service, state ga.State, data ga.EntityData) {
		}).
		Build()
	r.app.RegisterEntityListeners(entityListener)
	r.app.Start()

	return nil
}
