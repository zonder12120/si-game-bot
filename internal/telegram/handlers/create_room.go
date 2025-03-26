package handlers

import (
	"fmt"
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type CreateRoomHandler struct {
	botService  *service.Bot
	gameService *game.Service
}

func NewCreateRoomHandler(
	botService *service.Bot,
	gameService *game.Service,
) *CreateRoomHandler {
	return &CreateRoomHandler{
		botService:  botService,
		gameService: gameService,
	}
}

func (h *CreateRoomHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle create room")

	roomID := h.gameService.CreateRoom(tgID)

	s.UpdateState(state.OnNewRoundMenu)
	s.UpdateRoomID(roomID)

	_, err := h.botService.SendMsgAndGetInfoMdv2(
		tgID,
		fmt.Sprintf(render.MsgRoomCreated, roomID),
		keyboard.NewRoundMenu(),
	)
	if err != nil {
		return err
	}
	return nil
}
