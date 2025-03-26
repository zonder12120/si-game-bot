package game

import (
	"gopkg.in/telebot.v3"
)

func (s *Service) notifyAdmin(roomID string, msg string, keyboard *telebot.ReplyMarkup) error {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return err
	}

	adminTgID := foundRoom.Admin.TgID

	err = s.botService.SendMessage(adminTgID, msg, keyboard)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) notifyAllMembers(roomID string, msg string, keyboard *telebot.ReplyMarkup) (map[int64]int, error) {
	tgIDs, err := s.getAllMembersIDs(roomID)
	if err != nil {
		return nil, err
	}

	return s.notifyService.NotifyUsers(tgIDs, msg, keyboard)
}

func (s *Service) notifyAllPlayers(roomID string, msg string, keyboard *telebot.ReplyMarkup) (map[int64]int, error) {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	tgIDs := make([]int64, 0, len(foundRoom.Players))
	for _, player := range foundRoom.Players {
		tgIDs = append(tgIDs, player.TgID)
	}

	return s.notifyService.NotifyUsers(tgIDs, msg, keyboard)
}

func (s *Service) notifyAllActivePlayers(roomID string, msg string, keyboard *telebot.ReplyMarkup) error {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return err
	}

	activePlayerIDs := make([]int64, 0, len(foundRoom.Round.ActivePlayers))
	for _, player := range foundRoom.Round.ActivePlayers {
		activePlayerIDs = append(activePlayerIDs, player.TgID)
	}

	_, err = s.notifyService.NotifyUsers(activePlayerIDs, msg, keyboard)
	return err
}
