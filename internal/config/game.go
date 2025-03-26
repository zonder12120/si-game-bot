package config

import "time"

type Game struct {
	RoundTTL  time.Duration `env:"ROUND_TIMEOUT" envDefault:"15s"`
	AnswerTTL time.Duration `env:"ANSWER_TIMEOUT" envDefault:"10s"`
}
