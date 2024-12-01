package homeassistant

import (
	ga "saml.dev/gome-assistant"

	"github.com/angelokurtis/go-home-automations/pkg/app"
)

func NewServices(app *ga.App) *ga.Service {
	return app.GetService()
}

func NewLight(service *ga.Service) app.Light {
	return service.Light
}
