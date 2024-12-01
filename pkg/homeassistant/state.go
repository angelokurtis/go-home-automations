package homeassistant

import ga "saml.dev/gome-assistant"

func NewState(app *ga.App) ga.State {
	return app.GetState()
}
