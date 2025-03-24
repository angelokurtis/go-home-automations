package app

import (
	"context"
	"log/slog"
	"strings"

	"github.com/angelokurtis/go-otel/span"
	"github.com/lmittmann/tint"
	"github.com/samber/lo"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/errors"
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

	entities, err := r.state.ListEntities()
	if err != nil {
		return errors.WithStack(err)
	}

	entities = lo.Filter(entities, func(item ga.EntityState, index int) bool {
		return strings.HasPrefix(item.EntityID, "switch.") && item.EntityID != "switch.zigbee2mqtt_bridge_permit_join"
	})

	entityListeners := lo.Map(r.entityListeners, func(entityListener EntityListener, index int) ga.EntityListener {
		return ga.NewEntityListener().
			EntityIds(entityListener.EntityIds()...).
			Call(func(service *ga.Service, state ga.State, entity ga.EntityData) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.InfoContext(ctx, "Entity state changed",
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
	slog.InfoContext(context.Background(), "Registered entity listeners",
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
	slog.InfoContext(context.Background(), "Registered event listeners",
		slog.Int("count", len(eventListeners)),
	)

	dailyTasks := lo.Map(r.dailyTasks, func(dailyTask DailyTask, index int) ga.DailySchedule {
		return ga.NewDailySchedule().
			Call(func(service *ga.Service, state ga.State) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "")

				if err := dailyTask.Execute(ctx); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "")
			}).
			At(dailyTask.ScheduledTime()).
			Build()
	})
	r.app.RegisterSchedules(dailyTasks...)
	slog.InfoContext(context.Background(), "",
		slog.Int("count", len(dailyTasks)),
	)

	sunriseTasks := lo.Map(r.sunriseTasks, func(sunriseTask SunriseTask, index int) ga.DailySchedule {
		return ga.NewDailySchedule().
			Call(func(service *ga.Service, state ga.State) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "")

				if err := sunriseTask.Execute(ctx); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "")
			}).
			At(sunriseTask.SunriseOffset()).
			Build()
	})
	r.app.RegisterSchedules(sunriseTasks...)
	slog.InfoContext(context.Background(), "",
		slog.Int("count", len(sunriseTasks)),
	)

	sunsetTasks := lo.Map(r.sunsetTasks, func(sunsetTask SunsetTask, index int) ga.DailySchedule {
		return ga.NewDailySchedule().
			Call(func(service *ga.Service, state ga.State) {
				ctx, end := span.Start(ctx)
				defer end()

				slog.DebugContext(ctx, "")

				if err := sunsetTask.Execute(ctx); err != nil {
					_ = span.Error(ctx, err)
					slog.ErrorContext(ctx, "", tint.Err(err))
					return
				}

				slog.InfoContext(ctx, "")
			}).
			At(sunsetTask.SunsetOffset()).
			Build()
	})
	r.app.RegisterSchedules(sunsetTasks...)
	slog.InfoContext(context.Background(), "",
		slog.Int("count", len(sunsetTasks)),
	)

	r.app.Start()

	slog.InfoContext(ctx, "Application started")

	return nil
}
