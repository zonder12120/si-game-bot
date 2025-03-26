package handlers

import (
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
	"github.com/zonder12120/tg-quiz/internal/telegram/user/access"
)

type RoundMenuHandler struct {
	botService    *service.Bot
	accessChecker *access.Checker

	createRoundHandler *CreateRoundHandler
	endGameHandler     *EndGameHandler

	unknownCMDHandler *UnknownCmdHandler
}

func NewRoundMenuHandler(
	botService *service.Bot,
	accessChecker *access.Checker,

	createRoundHandler *CreateRoundHandler,
	endGameHandler *EndGameHandler,

	unknownCMDHandler *UnknownCmdHandler,
) *RoundMenuHandler {
	return &RoundMenuHandler{
		botService:    botService,
		accessChecker: accessChecker,

		createRoundHandler: createRoundHandler,
		endGameHandler:     endGameHandler,

		unknownCMDHandler: unknownCMDHandler,
	}
}

func (h *RoundMenuHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnNewRoundMenu
}

func (h *RoundMenuHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle new round menu")

	switch c.Text() {
	case render.BtnNewRound:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push new round")
		if !h.accessChecker.Check(c, s, access.CmdNewRound) {
			return h.botService.SendMessage(tgID, render.ErrMessageTextForbidden, nil)
		}
		return h.createRoundHandler.handleEnterPoints(c, s)
	case render.BtnEndGame:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push end game")
		if !h.accessChecker.Check(c, s, access.CmdEndGame) {
			return h.botService.SendMessage(tgID, render.ErrMessageTextForbidden, nil)
		}
		return h.endGameHandler.handleConfirm(c, s)
	default:
		return h.unknownCMDHandler.Handle(c, s)
	}
}
