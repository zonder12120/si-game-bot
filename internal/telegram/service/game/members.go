package game

import (
	"fmt"

	"github.com/zonder12120/tg-quiz/internal/room"
)

func (s *Service) getAllMembersIDs(roomID string) ([]int64, error) {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	if foundRoom.Admin == nil {
		return nil, fmt.Errorf("room has no admin")
	}
	if foundRoom.Players == nil {
		return nil, fmt.Errorf("room has no players")
	}

	tgIDs := make([]int64, 0, len(foundRoom.Players)+1)
	tgIDs = append(tgIDs, foundRoom.Admin.TgID)

	for _, player := range foundRoom.Players {
		if player != nil {
			tgIDs = append(tgIDs, player.TgID)
		}
	}
	return tgIDs, nil
}

func (s *Service) GetPlayer(roomID string, tgID int64) (*room.Member, error) {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	if foundRoom.Players == nil {
		return nil, fmt.Errorf("room has no players")
	}

	player, ok := foundRoom.Players[tgID]
	if !ok {
		return nil, fmt.Errorf("player not found")
	}

	return player, nil
}
