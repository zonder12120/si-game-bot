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

type LeaveRoomHandler struct {
	gameService *game.Service
	botService  *service.Bot
}

func NewLeaveRoomHandler(
	gameService *game.Service,
	botService *service.Bot,
) *LeaveRoomHandler {
	return &LeaveRoomHandler{
		gameService: gameService,
		botService:  botService,
	}
}

func (h *LeaveRoomHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnLeavingRoom
}

func (h *LeaveRoomHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle leave room")

	switch c.Text() {
	case render.BtnConfirm:
		err := h.handleLeave(c, s)
		if err != nil {
			return err
		}
	case render.BtnCancel:
		s.UpdateState(state.OnWaitingNewRound)
		err := h.botService.SendMessage(
			tgID,
			render.MsgWaitAnswering,
			keyboard.Leave(),
		)
		if err != nil {
			return err
		}
	default:
		return h.botService.SendMessage(tgID, render.MsgUnknownCommand, nil)
	}
	return nil
}

func (h *LeaveRoomHandler) handleConfirmLeave(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle confirm leave room")

	s.UpdateState(state.OnLeavingRoom)
	err := h.botService.SendMessage(
		tgID,
		render.MsgAreYouSure,
		keyboard.ConfirmationMenu(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (h *LeaveRoomHandler) handleLeave(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	err := h.gameService.LeaveRoom(s.GetRoomID(), tgID)
	if err != nil {
		return err
	}

	s.UpdateState(state.Idle)
	err = h.botService.SendMessage(
		tgID,
		render.MsgLeaveGame,
		keyboard.MainMenu(),
	)
	if err != nil {
		return err
	}
	s.UpdateRoomID("")

	return nil
}
