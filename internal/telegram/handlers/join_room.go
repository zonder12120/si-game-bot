package handlers

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type JoinRoomHandler struct {
	botService  *service.Bot
	gameService *game.Service
}

func NewJoinRoomHandler(
	botService *service.Bot,
	gameService *game.Service,
) *JoinRoomHandler {
	return &JoinRoomHandler{
		botService:  botService,
		gameService: gameService,
	}
}

func (h *JoinRoomHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnEnteredRoomID
}

func (h *JoinRoomHandler) Handle(c telebot.Context, s *state.UserSession) error {
	user := c.Sender()
	tgID := user.ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle join room")

	fullName := fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	if user.LastName == "" {
		fullName = user.FirstName
	}

	roomID := strings.TrimSpace(c.Text())
	if roomID == "" {
		return fmt.Errorf("empty roomID")
	}

	err := h.gameService.JoinRoom(roomID, tgID, fullName)
	if err != nil {
		s.UpdateState(state.Idle)
		sendErr := h.botService.SendMessage(
			tgID,
			render.MsgCantJoin,
			nil,
		)
		if sendErr != nil {
			return sendErr
		}
		return err
	}

	err = h.botService.SendMessage(
		tgID,
		render.MsgJoinedInTheRoom,
		keyboard.Leave(),
	)
	if err != nil {
		return h.botService.SendMessage(
			tgID,
			fmt.Sprintf(render.ErrUnexpected, err),
			keyboard.Leave(),
		)
	}

	s.UpdateState(state.OnWaitingNewRound)
	s.UpdateRoomID(roomID)

	return nil
}

func (h *JoinRoomHandler) handleEnterRoomID(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle ender room id")

	s.UpdateState(state.OnEnteredRoomID)

	return h.botService.SendMessage(
		tgID,
		render.MsgEnterRoomID,
		nil,
	)
}
