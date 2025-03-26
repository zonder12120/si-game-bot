package config

import "time"

type RoomManager struct {
	TTL time.Duration `env:"ROOM_TTL" envDefault:"30m"`
}
