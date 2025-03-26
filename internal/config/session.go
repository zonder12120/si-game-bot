package config

import "time"

type Session struct {
	TTL time.Duration `env:"SESSION_TTL" envDefault:"30m"`
}
