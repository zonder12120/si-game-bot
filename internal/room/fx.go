package room

import "go.uber.org/fx"

var Module = fx.Module(
	"room_manager",
	fx.Provide(NewManager),
)
