package homeassistant

import ga "saml.dev/gome-assistant"

func NewServices(app *ga.App) *ga.Service {
	return app.GetService()
}
