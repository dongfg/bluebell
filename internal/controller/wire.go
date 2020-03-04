package controller

import "github.com/google/wire"

// WireSet of controller
var WireSet = wire.NewSet(wire.Struct(new(Options), "*"), New)
