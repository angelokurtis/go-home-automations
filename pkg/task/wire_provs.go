package task

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewEveningLightsOn,
)
