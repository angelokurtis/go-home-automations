package action

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewAllSwitchesOff,
	NewPowerControl,
)
