package handlers

import (
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
)

type StartHandler struct {
	botService *service.Bot
}

func NewStartHandler(botService *service.Bot) *StartHandler {
	return &StartHandler{
		botService: botService,
	}
}

func (h *StartHandler) Handle(c telebot.Context) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle start")

	userSession := h.botService.SessionsManager.GetSession(tgID)
	userSession.Reset()

	return h.botService.SendMessage(
		tgID,
		render.MsgStartCommand,
		keyboard.MainMenu(),
	)
}
