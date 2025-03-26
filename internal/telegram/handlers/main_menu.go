package handlers

import (
	"strconv"

	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type MainMenuHandler struct {
	botService *service.Bot

	createRoomHandler *CreateRoomHandler
	joinRoomHandler   *JoinRoomHandler

	unknownCMDHandler *UnknownCmdHandler
}

func NewMainMenuHandler(
	botService *service.Bot,

	joinRoomHandler *JoinRoomHandler,
	createRoomHandler *CreateRoomHandler,

	unknownCMDHandler *UnknownCmdHandler,

) *MainMenuHandler {
	return &MainMenuHandler{
		botService: botService,

		joinRoomHandler:   joinRoomHandler,
		createRoomHandler: createRoomHandler,

		unknownCMDHandler: unknownCMDHandler,
	}
}

func (h *MainMenuHandler) CanHandle(currentState state.State) bool {
	return currentState == state.Idle
}

func (h *MainMenuHandler) Handle(c telebot.Context, s *state.UserSession) error {
	tgID := c.Sender().ID

	h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("handle main menu")

	switch c.Text() {
	case render.BtnJoinRoom:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push join room")
		return h.joinRoomHandler.handleEnterRoomID(c, s)
	case render.BtnCreateRoom:
		h.botService.Log.Debug().Str("user", strconv.Itoa(int(tgID))).Msg("push create room")
		return h.createRoomHandler.Handle(c, s)

	default:
		return h.unknownCMDHandler.Handle(c, s)
	}
}
