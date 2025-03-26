package internal

import (
	"go.uber.org/fx"

	"github.com/zonder12120/tg-quiz/internal/config"
	"github.com/zonder12120/tg-quiz/internal/logger"
	"github.com/zonder12120/tg-quiz/internal/room"
	"github.com/zonder12120/tg-quiz/internal/telegram"
)

func DeclareAppOpts() []fx.Option {
	return []fx.Option{
		fx.Provide(
			config.Parse,
		),

		room.Module,
		logger.Module,
		telegram.Module,
	}
}

func NewApp() *fx.App {
	return fx.New(DeclareAppOpts()...)
}
