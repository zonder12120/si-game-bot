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

func (h *LeaveRoomHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle leave room")

	err := h.gameService.LeaveRoom(s.GetRoomID(), tgID)
	if err != nil {
		return err
	}

	err = h.botService.SendMessage(
		tgID,
		render.MsgLeaveGame,
		keyboard.MainMenu(),
	)
	if err != nil {
		return err
	}

	s.UpdateState(state.Idle)
	s.UpdateRoomID("")

	return nil
}
