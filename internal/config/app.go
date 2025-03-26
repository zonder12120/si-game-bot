package config

import (
	"go.uber.org/fx"
)

type App struct {
	fx.Out

	Game        Game
	RoomManager RoomManager
	Session     Session
	Telegram    Telegram
	Logging     Logging
}
