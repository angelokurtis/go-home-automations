package homeassistant

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewApp,
	NewConfigFromEnv,
	NewServices,
	NewLight,
)
