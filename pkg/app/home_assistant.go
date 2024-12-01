package app

import ga "saml.dev/gome-assistant"

//go:generate go run -mod=mod github.com/vburenin/ifacemaker@latest -f "$GOMODCACHE/saml.dev/gome-assistant@v0.2.3/app.go" -s "App" -i "HomeAssistant" -p "app" -y "HomeAssistant defines methods for integrating with a Home Assistant instance" -o "home_assistant_gen.go"

// DailySchedule represents a task scheduled to run at a specific time each day.
type DailySchedule = ga.DailySchedule

// Interval represents a recurring task with a specified time interval.
type Interval = ga.Interval

// EntityListener listens for changes in the state of specific entities.
type EntityListener = ga.EntityListener

// EventListener listens for specific events on the Home Assistant event bus.
type EventListener = ga.EventListener

// Service provides core functionalities for managing Home Assistant integrations.
type Service = ga.Service

// State represents the current state of the Home Assistant instance, including entity states.
type State = ga.State
