package app

import ga "saml.dev/gome-assistant"

//go:generate go run -mod=mod github.com/vburenin/ifacemaker@latest -f "$GOMODCACHE/saml.dev/gome-assistant@v0.2.3/app.go" -s "App" -i "HomeAssistant" -p "app" -y "HomeAssistant defines methods for integrating with a Home Assistant instance" -o "home_assistant_gen.go"

type (
	DailySchedule  ga.DailySchedule
	Interval       ga.Interval
	EntityListener ga.EntityListener
	EventListener  ga.EventListener
	Service        ga.Service
	State          ga.State
)
