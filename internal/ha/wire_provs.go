package ha

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	LoadConfig,
	NewApp,
	NewService,
	NewState,
)
