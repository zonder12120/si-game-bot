package telegram

import (
	"context"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/zonder12120/tg-quiz/internal/telegram/bot"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/notify"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
	"github.com/zonder12120/tg-quiz/internal/telegram/user/access"
)

var Module = fx.Module("telegram_bot",
	fx.Provide(
		bot.NewBot,
		bot.NewWorker,
		service.NewBot,
		state.NewManager,
		state.NewDispatcher,
		access.NewAccessChecker,
		notify.NewService,
		game.NewService,
		NewHandlerRegister,
	),

	fx.Invoke(
		func(register *HandlerRegister) error {
			return register.RegisterHandlers()
		},
		registerWorkerLifecycle,
	),
)

func registerWorkerLifecycle(lifecycle fx.Lifecycle, worker *bot.Worker, log *zerolog.Logger) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info().Msg("Starting telegram worker")

			botCtx, cancel := context.WithCancel(context.Background())

			go func() {
				if err := worker.Run(botCtx); err != nil {
					log.Error().Err(err).Msg("Telegram worker failed")
				}
			}()

			lifecycle.Append(fx.Hook{
				OnStop: func(ctx context.Context) error {
					log.Info().Msg("Stopping telegram worker")
					cancel()
					worker.Stop()
					return nil
				},
			})

			return nil
		},
	})
}
