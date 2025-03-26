package notify

import (
	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/service"
)

type Service struct {
	botService *service.Bot
}

func NewService(botService *service.Bot) *Service {
	return &Service{
		botService: botService,
	}
}

func (s *Service) NotifyUsers(tgIDs []int64, msg string, keyboard *telebot.ReplyMarkup) (map[int64]int, error) {
	messageIDs := make(map[int64]int)

	for _, tgID := range tgIDs {
		message, err := s.botService.SendMsgAndGetInfo(tgID, msg, keyboard)
		if err != nil {
			return nil, err
		}
		messageIDs[tgID] = message.ID
	}

	return messageIDs, nil
}

func (s *Service) UpdateMessages(messageIDs map[int64]int, text string) {
	for tgID, msgID := range messageIDs {
		_, _ = s.botService.Bot.Edit(&telebot.Message{
			ID:   msgID,
			Chat: &telebot.Chat{ID: tgID},
		}, text)
	}
}
