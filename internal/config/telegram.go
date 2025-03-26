package config

type Telegram struct {
	Token string `env:"TELEGRAM_BOT_TOKEN,required"`
}
