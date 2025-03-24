package main

import (
	"github.com/google/wire"
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/internal/ha"
	"github.com/angelokurtis/go-home-automations/internal/term"
	"github.com/angelokurtis/go-home-automations/pkg/app"
	"github.com/angelokurtis/go-home-automations/pkg/entity"
	"github.com/angelokurtis/go-home-automations/pkg/task"
)

//nolint:unused // This function is used during compile-time to generate code for dependency injection
var providers = wire.NewSet(
	app.NewRunner,
	ha.Providers,
	term.Providers,
	wire.Bind(new(Runner), new(*app.Runner)),
	wire.Bind(new(term.Renderer), new(*term.MarkdownRenderer)),

	entityListeners,
	eventListeners,
	dailyTasks,
	sunsetTasks,
	sunriseTasks,
)

func entityListeners(service *ga.Service, state ga.State) []app.EntityListener {
	return []app.EntityListener{
		// arandelas de dentro
		entity.NewSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_6x_da_copa_l1",
			"switch.interruptor_6x_da_copa_l2",
			"switch.interruptor_3x_da_entrada_left",
		),
		// luz da sala
		entity.NewSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_3x_da_entrada_center",
			"switch.interruptor_6x_da_entrada_l1",
		),
		// luz da entrada
		entity.NewSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_3x_da_entrada_right",
			"switch.interruptor_6x_da_entrada_l2",
		),
		// luz da churrasqueira
		entity.NewSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_2x_da_area_gourmet_center",
			"switch.interruptor_6x_da_area_gourmet_l5",
			"switch.interruptor_6x_da_copa_l5",
		),
		// arandelas da piscina
		entity.NewSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_6x_da_area_gourmet_l1",
			"switch.interruptor_6x_da_area_gourmet_l6",
		),
		// arandela da lateral
		entity.NewSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_2x_da_area_gourmet_left",
			"switch.interruptor_6x_da_copa_l6",
		),
		// luzes do jardim
		entity.NewSynchronizedSwitchesListener(
			service,
			state,
			"switch.interruptor_6x_da_entrada_l3",
			"switch.interruptor_6x_da_entrada_l4",
		),
	}
}

func eventListeners() []app.EventListener {
	return []app.EventListener{}
}

func dailyTasks(service *ga.Service, state ga.State) []app.DailyTask {
	return []app.DailyTask{
		task.NewAllSwitchesOff(service, state),
	}
}

func sunsetTasks() []app.SunsetTask {
	return []app.SunsetTask{}
}

func sunriseTasks() []app.SunriseTask {
	return []app.SunriseTask{}
}
