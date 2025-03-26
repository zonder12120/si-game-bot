package access

import (
	"fmt"

	"github.com/rs/zerolog"
	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/room"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type Checker struct {
	roomManager *room.Manager
	log         *zerolog.Logger
}

func NewAccessChecker(roomManager *room.Manager, log *zerolog.Logger) *Checker {
	return &Checker{
		roomManager: roomManager,
		log:         log,
	}
}

func (a *Checker) Check(c telebot.Context, s *state.UserSession, command Command) bool {
	tgID := c.Sender().ID

	player, err := a.roomManager.GetMember(s.GetRoomID(), tgID)
	if err != nil {
		return a.sendError(c, fmt.Sprint(err))
	}

	if !HasAccess(player.Role, command) {
		return a.sendError(c, render.ErrMessageTextForbidden)
	}

	return true
}

func (a *Checker) sendError(c telebot.Context, msg string) bool {
	err := c.Send(msg)
	if err != nil {
		a.log.Error().Err(err).Msg("failed to send error message")
	}
	return false
}
