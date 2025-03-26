package state

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Session struct {
	user   sync.Map
	ttl    time.Duration
	logger *zerolog.Logger
}

type UserSession struct {
	mu           sync.RWMutex
	TelegramID   int64
	CurrentState State
	RoomID       string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	log          *zerolog.Logger
}

func NewSession(ttl time.Duration, logger *zerolog.Logger) *Session {
	s := &Session{
		ttl:    ttl,
		logger: logger,
	}
	go s.cleanupWorker()
	return s
}

func NewUserSession(telegramID int64, log *zerolog.Logger) *UserSession {
	return &UserSession{
		TelegramID:   telegramID,
		CurrentState: Idle,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		log:          log,
	}
}

func (s *Session) get(telegramID int64) *UserSession {
	if session, ok := s.user.Load(telegramID); ok {
		return session.(*UserSession)
	}

	newSession := NewUserSession(telegramID, s.logger)
	session, _ := s.user.LoadOrStore(telegramID, newSession)
	return session.(*UserSession)
}

func (s *Session) cleanupWorker() {
	ticker := time.NewTicker(s.ttl)
	defer ticker.Stop()

	for range ticker.C {
		s.user.Range(func(key, value interface{}) bool {
			session := value.(*UserSession)
			session.mu.RLock()
			updatedAt := session.UpdatedAt
			session.mu.RUnlock()

			if time.Since(updatedAt) > s.ttl {
				s.user.Delete(key)
				s.logger.Info().Int64("telegramID", key.(int64)).Msg("Session cleaned up")
			}
			return true
		})
	}
}

func (s *UserSession) UpdateState(newState State) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Debug().Msg(fmt.Sprintf("Updating state from %v to %v", s.CurrentState, newState))
	s.CurrentState = newState
	s.UpdatedAt = time.Now()
}

func (s *UserSession) UpdateRoomID(newRoomID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Debug().Msg(fmt.Sprintf("Updating room from %v to %v", s.RoomID, newRoomID))
	s.RoomID = newRoomID
	s.UpdatedAt = time.Now()
}

func (s *UserSession) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	log.Debug().Msg(fmt.Sprintf("Resetting session for %v", s.TelegramID))
	s.CurrentState = Idle
	s.RoomID = ""
	s.UpdatedAt = time.Now()
}

func (s *UserSession) GetState() (State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Debug().Msg(fmt.Sprintf("Getting state from user session tgID: %v", s.TelegramID))
	return s.CurrentState, nil
}

func (s *UserSession) GetRoomID() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	log.Debug().Msg(fmt.Sprintf("Getting room from user session tgID: %v", s.TelegramID))
	return s.RoomID
}
