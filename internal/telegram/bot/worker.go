package bot

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"gopkg.in/telebot.v3"
)

type Worker struct {
	bot      *telebot.Bot
	log      *zerolog.Logger
	stopOnce sync.Once
	stopChan chan struct{}
}

func NewWorker(bot *telebot.Bot, log *zerolog.Logger) *Worker {
	return &Worker{
		bot:      bot,
		log:      log,
		stopChan: make(chan struct{}),
	}
}

func (w *Worker) Run(ctx context.Context) error {
	w.log.Info().Msg("Starting telegram bot worker")

	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		w.bot.Start()
		close(done)
	}()

	select {
	case <-ctx.Done():
		w.log.Info().Msg("Received shutdown signal, stopping bot")
		w.Stop()

		select {
		case <-done:
			w.log.Info().Msg("Bot stopped gracefully")
		case <-time.After(5 * time.Second):
			w.log.Warn().Msg("Force stopping bot")
		}
		return ctx.Err()

	case <-done:
		w.log.Info().Msg("Bot stopped internally")
		return nil
	}
}

func (w *Worker) Stop() {
	w.stopOnce.Do(func() {
		w.log.Info().Msg("Initiating bot stop")
		close(w.stopChan)
		if w.bot != nil {
			w.log.Debug().Msg("Calling bot.Stop()")
			w.bot.Stop()
			w.log.Info().Msg("Bot stopped successfully")
		}
	})
}
