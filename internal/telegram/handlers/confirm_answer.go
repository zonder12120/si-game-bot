package handlers

import (
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
	"github.com/zonder12120/tg-quiz/internal/telegram/user/access"
)

// state.OnConfirmAnswer

type ConfirmAnswerHandler struct {
	botService    *service.Bot
	gameService   *game.Service
	accessChecker *access.Checker
}

func NewConfirmAnswerHandler(
	botService *service.Bot,
	gameService *game.Service,
	accessChecker *access.Checker,
) *ConfirmAnswerHandler {
	return &ConfirmAnswerHandler{
		botService:    botService,
		gameService:   gameService,
		accessChecker: accessChecker,
	}
}

func (h *ConfirmAnswerHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnConfirmAnswer
}

func (h *ConfirmAnswerHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle confirm answer")

	switch c.Text() {
	case render.BtnYes:
		if !h.accessChecker.Check(c, s, access.CmdConfirmAnswer) {
			return h.botService.SendMessage(tgID, render.ErrMessageTextForbidden, nil)
		}
		err := h.gameService.ResultAnswer(s.GetRoomID(), true)
		if err != nil {
			return err
		}
	case render.BtnNo:
		if !h.accessChecker.Check(c, s, access.CmdConfirmAnswer) {
			return h.botService.SendMessage(tgID, render.ErrMessageTextForbidden, nil)
		}
		err := h.gameService.ResultAnswer(s.GetRoomID(), false)
		if err != nil {
			return err
		}
	default:
		return h.botService.SendMessage(tgID, render.MsgUnknownCommand, nil)
	}
	return nil
}
