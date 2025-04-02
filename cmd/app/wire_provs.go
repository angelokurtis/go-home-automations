package main

import (
	"github.com/google/wire"

	"github.com/angelokurtis/go-home-automations/internal/ha"
	"github.com/angelokurtis/go-home-automations/internal/term"
	"github.com/angelokurtis/go-home-automations/pkg/action"
	"github.com/angelokurtis/go-home-automations/pkg/app"
	"github.com/angelokurtis/go-home-automations/pkg/entity"
	"github.com/angelokurtis/go-home-automations/pkg/event"
	"github.com/angelokurtis/go-home-automations/pkg/task"
)

//nolint:unused // This function is used during compile-time to generate code for dependency injection
var providers = wire.NewSet(
	app.NewRunner,
	ha.Providers,
	term.Providers,
	wire.Bind(new(app.PowerControl), new(*action.PowerControl)),
	wire.Bind(new(Runner), new(*app.Runner)),
	wire.Bind(new(term.Renderer), new(*term.MarkdownRenderer)),
	action.Providers,
	entity.Providers,
	event.Providers,
	task.Providers,
	entityListeners,
	eventListeners,
	dailyTasks,
	sunsetTasks,
	sunriseTasks,
)

func entityListeners(synchronizedSwitches entity.SynchronizedSwitchesListeners) []app.EntityListener {
	var listeners []app.EntityListener
	for _, listener := range synchronizedSwitches {
		listeners = append(listeners, listener)
	}

	return listeners
}

func eventListeners() []app.EventListener {
	return []app.EventListener{}
}

func dailyTasks(allSwitchesOff *action.AllSwitchesOff) []app.DailyTask {
	return []app.DailyTask{
		task.NewDaily(allSwitchesOff, "23:30", "Apagar Todas os Interruptores"),
	}
}

func sunsetTasks(eveningLights *task.EveningLightsOn) []app.SunsetTask {
	return []app.SunsetTask{
		eveningLights,
	}
}

func sunriseTasks() []app.SunriseTask {
	return []app.SunriseTask{}
}
