package telegram

import (
	"go.uber.org/fx"
	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/handlers"
	"github.com/zonder12120/tg-quiz/internal/telegram/middleware"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/game"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
	"github.com/zonder12120/tg-quiz/internal/telegram/user/access"
)

type HandlerRegister struct {
	botService      *service.Bot
	stateDispatcher *state.Dispatcher
	accessChecker   *access.Checker
	gameService     *game.Service
}

type NewHandlerRegisterParams struct {
	fx.In

	BotService      *service.Bot
	StateDispatcher *state.Dispatcher
	AccessChecker   *access.Checker
	GameService     *game.Service
}

func NewHandlerRegister(params NewHandlerRegisterParams) *HandlerRegister {
	return &HandlerRegister{
		botService:      params.BotService,
		stateDispatcher: params.StateDispatcher,
		accessChecker:   params.AccessChecker,
		gameService:     params.GameService,
	}
}

func (h *HandlerRegister) RegisterHandlers() error {
	// init handlers
	unknownCMDHandler := handlers.NewUnknownCmdHandler(h.botService)
	startHandler := handlers.NewStartHandler(h.botService)
	cancelHandler := handlers.NewCancelHandler(h.botService)

	createRoomHandler := handlers.NewCreateRoomHandler(h.botService, h.gameService)
	endGameHandler := handlers.NewEndGameHandler(h.gameService, h.botService)
	createRoundHandler := handlers.NewCreateRoundHandler(h.botService, h.gameService)
	startRoundHandler := handlers.NewStartRoundHandler(h.gameService, h.botService, h.accessChecker)
	adminRoundHandler := handlers.NewAdminRoundHandler(h.gameService, h.botService)
	confirmAnswerHandler := handlers.NewConfirmAnswerHandler(h.botService, h.gameService, h.accessChecker)

	joinRoomHandler := handlers.NewJoinRoomHandler(h.botService, h.gameService)
	leaveRoomHandler := handlers.NewLeaveRoomHandler(h.gameService, h.botService)

	playingHandler := handlers.NewPlayingHandler(h.gameService, h.botService, h.accessChecker, leaveRoomHandler)

	newRoundMenuHandler := handlers.NewRoundMenuHandler(
		h.botService,
		h.accessChecker,
		createRoundHandler,
		endGameHandler,
		unknownCMDHandler,
	)

	mainMenuHandler := handlers.NewMainMenuHandler(
		h.botService,
		joinRoomHandler,
		createRoomHandler,
		unknownCMDHandler,
	)

	h.stateDispatcher.AddHandler(endGameHandler)
	h.stateDispatcher.AddHandler(createRoundHandler)
	h.stateDispatcher.AddHandler(newRoundMenuHandler)
	h.stateDispatcher.AddHandler(startRoundHandler)
	h.stateDispatcher.AddHandler(adminRoundHandler)
	h.stateDispatcher.AddHandler(confirmAnswerHandler)

	h.stateDispatcher.AddHandler(joinRoomHandler)
	h.stateDispatcher.AddHandler(leaveRoomHandler)

	h.stateDispatcher.AddHandler(playingHandler)

	h.stateDispatcher.AddHandler(mainMenuHandler)

	middlewares := []telebot.MiddlewareFunc{
		middleware.Error(h.botService.Log),
	}

	h.botService.Bot.Handle("/start", startHandler.Handle, middlewares...)
	h.botService.Bot.Handle("/reset", cancelHandler.Handle, middlewares...)

	// Handles text messages based on user state
	h.botService.Bot.Handle(telebot.OnText, func(c telebot.Context) error {
		return h.stateDispatcher.Process(c)
	}, middlewares...)

	return nil
}
