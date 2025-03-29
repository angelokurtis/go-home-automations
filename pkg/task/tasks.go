package task

import "context"

type Executor interface {
	Execute(ctx context.Context) error
}

type Daily struct {
	ex            Executor
	scheduledTime string
	name          string
}

func NewDaily(ex Executor, scheduledTime, name string) *Daily {
	return &Daily{ex: ex, scheduledTime: scheduledTime, name: name}
}

func (d *Daily) Execute(ctx context.Context) error {
	return d.ex.Execute(ctx)
}

func (d *Daily) ScheduledTime() string {
	return d.scheduledTime
}

func (d *Daily) Name() string {
	return d.name
}
