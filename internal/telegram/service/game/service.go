package game

import (
	"sync"

	"github.com/rs/zerolog"
	"go.uber.org/fx"

	"github.com/zonder12120/tg-quiz/internal/config"
	"github.com/zonder12120/tg-quiz/internal/room"
	"github.com/zonder12120/tg-quiz/internal/telegram/service"
	"github.com/zonder12120/tg-quiz/internal/telegram/service/notify"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type Service struct {
	roomManager    *room.Manager
	botService     *service.Bot
	notifyService  *notify.Service
	sessionManager *state.Manager
	gameConfig     config.Game
	log            *zerolog.Logger
	mu             sync.Mutex
	timers         map[string]map[TimerType]*Timer
}

type NewServiceParams struct {
	fx.In

	RoomManager    *room.Manager
	BotService     *service.Bot
	NotifyService  *notify.Service
	SessionManager *state.Manager
	GameConfig     config.Game
	Log            *zerolog.Logger
}

func NewService(params NewServiceParams) *Service {
	return &Service{
		roomManager:    params.RoomManager,
		botService:     params.BotService,
		notifyService:  params.NotifyService,
		sessionManager: params.SessionManager,
		gameConfig:     params.GameConfig,
		log:            params.Log,
		timers:         make(map[string]map[TimerType]*Timer),
	}
}
