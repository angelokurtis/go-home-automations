package app

import (
	"github.com/google/wire"
)

var Providers = wire.NewSet(
	NewRunner,
)
