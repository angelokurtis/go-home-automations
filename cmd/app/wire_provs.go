package main

import (
	"github.com/google/wire"
	ga "saml.dev/gome-assistant"

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
	action.Providers,
	app.NewRunner,
	entity.Providers,
	event.Providers,
	ha.Providers,
	task.Providers,
	term.Providers,
	wire.Bind(new(app.EntityReader), new(ga.State)),
	wire.Bind(new(app.PowerControl), new(*action.PowerControl)),
	wire.Bind(new(Runner), new(*app.Runner)),
	wire.Bind(new(term.Renderer), new(*term.MarkdownRenderer)),
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
		task.NewDaily(allSwitchesOff, "22:50", "Apagar Todas os Interruptores"),
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
