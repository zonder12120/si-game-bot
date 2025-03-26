package handlers

import (
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type EndGameHandler struct {
	gameService *game.Service
	botService  *service.Bot
}

func NewEndGameHandler(
	gameService *game.Service,
	botService *service.Bot,
) *EndGameHandler {
	return &EndGameHandler{
		gameService: gameService,
		botService:  botService,
	}
}

func (h *EndGameHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnConfirmEndGame
}

func (h *EndGameHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle end game")

	switch c.Text() {
	case render.BtnConfirm:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push confirm button")
		err := h.gameService.EndGame(s.GetRoomID())
		if err != nil {
			return err
		}
		return nil
	case render.BtnCancel:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push cancel button")
		s.UpdateState(state.OnNewRoundMenu)
		err := h.botService.SendMessage(tgID, render.MsgCreateRound, keyboard.NewRoundMenu())
		if err != nil {
			return err
		}
	}
	return h.botService.SendMessage(tgID, render.MsgUnknownCommand, nil)
}

func (h *EndGameHandler) handleConfirm(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle confirm end game")

	err := h.botService.SendMessage(tgID, render.MsgAreYouSure, keyboard.ConfirmationMenu())
	if err != nil {
		return err
	}

	s.UpdateState(state.OnConfirmEndGame)

	return nil
}
