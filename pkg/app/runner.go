package app

import (
	"context"
	"log/slog"

	"github.com/angelokurtis/go-otel/span"
	"github.com/lmittmann/tint"
	"github.com/samber/lo"
	ga "saml.dev/gome-assistant"
)

type Runner struct {
	app             *ga.App
	entityListeners []EntityListener
	eventListeners  []EventListener
	dailyTasks      []DailyTask
	sunsetTasks     []SunsetTask
	sunriseTasks    []SunriseTask

	svc   *ga.Service
	state ga.State
}

func NewRunner(app *ga.App, entityListeners []EntityListener, eventListeners []EventListener, dailyTasks []DailyTask, sunsetTasks []SunsetTask, sunriseTasks []SunriseTask, svc *ga.Service, state ga.State) *Runner {
	return &Runner{app: app, entityListeners: entityListeners, eventListeners: eventListeners, dailyTasks: dailyTasks, sunsetTasks: sunsetTasks, sunriseTasks: sunriseTasks, svc: svc, state: state}
}

func (r *Runner) Run(ctx context.Context) error {
	slog.InfoContext(ctx, "Starting application")

	entityListeners := lo.Map(r.entityListeners, func(entityListener EntityListener, index int) ga.EntityListener {
		return ga.NewEntityListener().
			EntityIds(entityListener.EntityIds()...).
			Call(func(service *ga.Service, state ga.State, entity ga.EntityData) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "Entity state changed",
					slog.String("entity_id", entity.TriggerEntityId),
					slog.String("new_state", entity.ToState),
				)

				if err := entityListener.OnChange(ctx, entity); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "Failed to process entity change", tint.Err(err))
					return
				}
			}).
			Build()
	})
	r.app.RegisterEntityListeners(entityListeners...)
	slog.InfoContext(ctx, "Registered entity listeners",
		slog.Int("count", len(entityListeners)),
	)

	eventListeners := lo.Map(r.eventListeners, func(eventListener EventListener, index int) ga.EventListener {
		return ga.NewEventListener().
			EventTypes(eventListener.EventTypes()...).
			Call(func(service *ga.Service, state ga.State, event ga.EventData) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "Event received",
					slog.String("event_type", event.Type),
					slog.String("raw", string(event.RawEventJSON)),
				)

				if err := eventListener.OnEvent(ctx, event); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "Failed to process event", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "Event processed")
			}).
			Build()
	})
	r.app.RegisterEventListeners(eventListeners...)
	slog.InfoContext(ctx, "Registered event listeners",
		slog.Int("count", len(eventListeners)),
	)

	dailyTasks := lo.Map(r.dailyTasks, func(dailyTask DailyTask, index int) ga.DailySchedule {
		return ga.NewDailySchedule().
			Call(func(service *ga.Service, state ga.State) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "Executing daily task",
					slog.String("task_name", dailyTask.Name()),
				)

				if err := dailyTask.Execute(ctx); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "Failed to execute daily task", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "Daily task executed")
			}).
			At(dailyTask.ScheduledTime()).
			Build()
	})
	r.app.RegisterSchedules(dailyTasks...)
	slog.InfoContext(ctx, "Registered daily tasks",
		slog.Int("count", len(dailyTasks)),
	)

	sunriseTasks := lo.Map(r.sunriseTasks, func(sunriseTask SunriseTask, index int) ga.DailySchedule {
		return ga.NewDailySchedule().
			Call(func(service *ga.Service, state ga.State) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "Executing sunrise task",
					slog.String("task_name", sunriseTask.Name()),
				)

				if err := sunriseTask.Execute(ctx); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "Failed to execute sunrise task", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "Sunrise task executed")
			}).
			Sunrise(ga.DurationString(sunriseTask.SunriseOffset())).
			Build()
	})
	r.app.RegisterSchedules(sunriseTasks...)
	slog.InfoContext(ctx, "Registered sunrise tasks",
		slog.Int("count", len(sunriseTasks)),
	)

	sunsetTasks := lo.Map(r.sunsetTasks, func(sunsetTask SunsetTask, index int) ga.DailySchedule {
		return ga.NewDailySchedule().
			Call(func(service *ga.Service, state ga.State) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "Executing sunset task",
					slog.String("task_name", sunsetTask.Name()),
				)

				if err := sunsetTask.Execute(ctx); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "Failed to execute sunset task", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "Sunset task executed")
			}).
			Sunset(ga.DurationString(sunsetTask.SunsetOffset())).
			Build()
	})
	r.app.RegisterSchedules(sunsetTasks...)
	slog.InfoContext(ctx, "Registered sunset tasks",
		slog.Int("count", len(sunsetTasks)),
	)

	r.app.Start()

	slog.InfoContext(ctx, "Application started")

	return nil
}
