package app

type SwitchController struct {
	light Light
}

func NewSwitchController(light Light) *SwitchController {
	return &SwitchController{light: light}
}

func (sc *SwitchController) OnStateChanged(event *StateChangeEvent) error {
	return nil
}
