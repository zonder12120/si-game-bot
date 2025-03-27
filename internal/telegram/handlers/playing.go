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

var playingStates = map[state.State]struct{}{
	state.OnWaitingNewRound: {},
	state.OnPlayingRound:    {},
}

type PlayingHandler struct {
	gameService   *game.Service
	botService    *service.Bot
	accessChecker *access.Checker

	leaveRoomHandler *LeaveRoomHandler
}

func NewPlayingHandler(
	gameService *game.Service,
	botService *service.Bot,
	accessChecker *access.Checker,

	leaveRoomHandler *LeaveRoomHandler,
) *PlayingHandler {
	return &PlayingHandler{
		gameService:   gameService,
		botService:    botService,
		accessChecker: accessChecker,

		leaveRoomHandler: leaveRoomHandler,
	}
}

func (h *PlayingHandler) CanHandle(currentState state.State) bool {
	_, ok := playingStates[currentState]
	return ok
}

func (h *PlayingHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle playing")

	switch s.CurrentState {
	case state.OnWaitingNewRound:
		return h.handleLeave(c, s)
	case state.OnPlayingRound:
		return h.handlePlaying(c, s)
	}
	return nil
}

func (h *PlayingHandler) handleLeave(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID

	if c.Text() == render.BtnLeave {
		err := h.leaveRoomHandler.handleConfirmLeave(c, s)
		if err != nil {
			return err
		}
	} else {
		return h.botService.SendMessage(tgID, render.MsgUnknownCommand, nil)
	}
	return nil
}

func (h *PlayingHandler) handlePlaying(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID

	switch c.Text() {
	case render.BtnLeave:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push leave")

		if !h.accessChecker.Check(c, s, access.CmdLeave) {
			return h.botService.SendMessage(tgID, render.ErrMessageTextForbidden, nil)
		}
		return h.handleLeave(c, s)
	case render.BtnAnswer:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push answer")

		if !h.accessChecker.Check(c, s, access.CmdAnswer) {
			return h.botService.SendMessage(tgID, render.ErrMessageTextForbidden, nil)
		}
		return h.handleAnswering(c, s)
	}
	return h.botService.SendMessage(tgID, render.MsgUnknownCommand, nil)
}

func (h *PlayingHandler) handleAnswering(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID
	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle answering")

	err := h.gameService.StartAnswer(s.GetRoomID(), tgID)
	if err != nil {
		return err
	}

	return nil
}
