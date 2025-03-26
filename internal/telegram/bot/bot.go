package bot

import (
	"net/http"
	"time"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/config"
)

func NewBot(cfg config.Telegram) (*telebot.Bot, error) {
	return telebot.NewBot(telebot.Settings{
		Token:  cfg.Token,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Client: &http.Client{
			Timeout: 15 * time.Second,
		},
	})
}
