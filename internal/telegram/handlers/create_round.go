package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type CreateRoundHandler struct {
	botService  *service.Bot
	gameService *game.Service
}

func NewCreateRoundHandler(
	botService *service.Bot,
	gameService *game.Service,
) *CreateRoundHandler {
	return &CreateRoundHandler{
		botService:  botService,
		gameService: gameService,
	}
}

func (h *CreateRoundHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnEnteredPoints
}

func (h *CreateRoundHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle create round")

	switch c.Text() {
	case render.BtnBack:
		s.UpdateState(state.OnNewRoundMenu)

		return h.botService.SendMessage(
			tgID,
			render.MsgCreateRound,
			keyboard.NewRoundMenu(),
		)
	default:
		return h.handleNewRound(c, s)
	}

}

func (h *CreateRoundHandler) handleEnterPoints(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle create room")

	s.UpdateState(state.OnEnteredPoints)

	return h.botService.SendMessage(
		tgID,
		render.MsgEnterPoints,
		keyboard.Back(),
	)
}

func (h *CreateRoundHandler) handleNewRound(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID

	strPoints := extractDigits(c.Text())
	if strPoints == "" {
		return fmt.Errorf("empty points")
	}

	points, err := strconv.Atoi(strPoints)
	if err != nil {
		return err
	}

	err = h.gameService.NewRound(s.GetRoomID(), points)
	if err != nil {
		return h.botService.SendMessage(tgID, fmt.Sprint(err), nil)
	}

	err = h.botService.SendMessage(tgID, render.MsgRoundCreated, keyboard.StartRoundMenu())
	if err != nil {
		return err
	}

	s.UpdateState(state.OnStartingRound)

	return nil
}

func extractDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}
