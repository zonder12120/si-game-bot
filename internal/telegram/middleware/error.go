package middleware

import (
	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"
)

func Error(log *zerolog.Logger) telebot.MiddlewareFunc {
	return func(next telebot.HandlerFunc) telebot.HandlerFunc {
		return func(c telebot.Context) error {
			err := next(c)
			if err != nil {
				log.Error().
					Err(err).
					Fields(struct {
						ChatID int64
						Text   string
					}{
						ChatID: c.Chat().ID,
						Text:   c.Text(),
					}).
					Msg("handler error")
				return c.Send("😔 Error: " + err.Error())
			}

			return nil
		}
	}
}
