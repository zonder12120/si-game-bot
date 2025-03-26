package state

import (
	"gopkg.in/telebot.v3"
)

type Handler interface {
	CanHandle(state State) bool
	Handle(c telebot.Context, s *UserSession) error
}

type Dispatcher struct {
	handlers []Handler
	manager  *Manager
}

func NewDispatcher(manager *Manager) *Dispatcher {
	return &Dispatcher{
		manager: manager,
	}
}

func (s *Dispatcher) AddHandler(handler Handler) {
	s.handlers = append(s.handlers, handler)
}

func (s *Dispatcher) Process(c telebot.Context) error {
	userSession := s.manager.GetSession(c.Sender().ID)

	userState, err := userSession.GetState()
	if err != nil {
		return err
	}

	for _, handler := range s.handlers {
		if handler.CanHandle(userState) {
			return handler.Handle(c, userSession)
		}
	}
	return nil
}
