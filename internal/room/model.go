package room

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog"

	"github.com/zonder12120/tg-quiz/internal/telegram/user"
)

type Manager struct {
	mu     sync.RWMutex
	rooms  sync.Map
	ttl    time.Duration
	log    *zerolog.Logger
	stopCh chan struct{}
	wg     sync.WaitGroup
}

type Room struct {
	mu        sync.RWMutex
	Admin     *Member
	Players   map[int64]*Member
	Round     *Round
	Status    Status
	updatedAt atomic.Value
}

type Member struct {
	TgID   int64
	Name   string
	Role   user.Role
	Points int
}

type Round struct {
	ActivePlayers   map[int64]*Member
	AnsweringPlayer int64
	Points          int
}

type Player struct{}

type Status uint8

const (
	waiting Status = iota
	playing
)

type RoomSnapshot struct {
	Admin   *Member
	Players map[int64]*Member
	Round   *RoundSnapshot
	Status  Status
}

type RoundSnapshot struct {
	ActivePlayers   map[int64]*Member
	AnsweringPlayer int64
	Points          int
}
