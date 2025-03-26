package handlers

import (
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
	"github.com/zonder12120/tg-quiz/internal/telegram/user/access"
)

type StartRoundHandler struct {
	gameService   *game.Service
	botService    *service.Bot
	accessChecker *access.Checker
}

func NewStartRoundHandler(
	gameService *game.Service,
	botService *service.Bot,
	accessChecker *access.Checker,
) *StartRoundHandler {
	return &StartRoundHandler{
		gameService:   gameService,
		botService:    botService,
		accessChecker: accessChecker,
	}
}

func (h *StartRoundHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnStartingRound
}

func (h *StartRoundHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle start round")

	switch c.Text() {
	case render.BtnRoundStart:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push start round")
		if !h.accessChecker.Check(c, s, access.CmdStartRound) {
			return h.botService.SendMessage(tgID, render.ErrMessageTextForbidden, nil)
		}
		s.UpdateState(state.OnAdminRound)

		err := h.botService.SendMessage(tgID, render.MsgRoundStarted, keyboard.EndRound())
		if err != nil {
			return err
		}

		err = h.gameService.StartRound(s.GetRoomID())
		if err != nil {
			return err
		}

	case render.BtnBack:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push back")
		s.UpdateState(state.OnNewRoundMenu)
		return h.botService.SendMessage(tgID, render.MsgCreateRound, keyboard.NewRoundMenu())

	default:
		return h.botService.SendMessage(tgID, render.MsgUnknownCommand, nil)
	}

	return nil
}
