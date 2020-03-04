package service

import "github.com/google/wire"

// WireSet of service
var WireSet = wire.NewSet(NewHealthService, NewTotpService, NewSeriesService)

// WireSet of service options
var OptionsSet = wire.NewSet(wire.Struct(new(HealthServiceOptions), "*"), wire.Struct(new(SeriesServiceOptions), "*"))
