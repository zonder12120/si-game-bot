package state

import (
	"github.com/rs/zerolog"

	"github.com/zonder12120/tg-quiz/internal/config"
)

type Manager struct {
	session *Session
	logger  *zerolog.Logger
}

func NewManager(session config.Session, logger *zerolog.Logger) *Manager {
	return &Manager{
		session: NewSession(session.TTL, logger),
	}
}

func (m *Manager) GetSession(userID int64) *UserSession {
	return m.session.get(userID)
}
