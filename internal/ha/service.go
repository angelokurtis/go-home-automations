package ha

import ga "saml.dev/gome-assistant"

func NewService(app *ga.App) *ga.Service {
	return app.GetService()
}
