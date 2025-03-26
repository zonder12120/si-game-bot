package handlers

import (
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
)

type CancelHandler struct {
	botService *service.Bot
}

func NewCancelHandler(botService *service.Bot) *CancelHandler {
	return &CancelHandler{
		botService: botService,
	}
}

func (h *CancelHandler) Handle(c telebot.Context) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle cancel")

	userSession := h.botService.SessionsManager.GetSession(tgID)
	userSession.Reset()

	return h.botService.SendMessage(
		tgID,
		render.MsgCancel,
		keyboard.MainMenu(),
	)
}
