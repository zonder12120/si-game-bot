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

type AdminRoundHandler struct {
	gameService *game.Service
	botService  *service.Bot
}

func NewAdminRoundHandler(
	gameService *game.Service,
	botService *service.Bot,
) *AdminRoundHandler {
	return &AdminRoundHandler{
		gameService: gameService,
		botService:  botService,
	}
}

func (h *AdminRoundHandler) CanHandle(currentState state.State) bool {
	return currentState == state.OnAdminRound
}

func (h *AdminRoundHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle admin round")

	if c.Text() == render.BtnEndRound {
		err := h.gameService.EndRound(s.GetRoomID())
		if err != nil {
			return err
		}

		s.UpdateState(state.OnNewRoundMenu)
		return h.botService.SendMessage(
			tgID,
			render.MsgCreateRound,
			keyboard.NewRoundMenu(),
		)
	}
	return nil
}
