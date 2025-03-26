package handlers

import (
	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type UnknownCmdHandler struct {
	botService *service.Bot
}

func NewUnknownCmdHandler(botService *service.Bot) *UnknownCmdHandler {
	return &UnknownCmdHandler{botService: botService}
}

func (h *UnknownCmdHandler) Handle(c telebot.Context, s *state.UserSession) error {
	userState, err := s.GetState()
	if err != nil {
		return err
	}
	h.botService.Log.
		Debug().
		Fields(struct {
			text  string
			state state.State
		}{
			text:  c.Text(),
			state: userState,
		}).
		Msg("unknown command")

	return h.botService.SendMessage(
		c.Sender().ID,
		render.MsgUnknownCommand,
		keyboard.Forbidden(),
	)
}
