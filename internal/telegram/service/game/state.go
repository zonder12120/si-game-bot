package game

import (
	"fmt"

	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

func (s *Service) updateAllPlayersState(roomID string, newState state.State) error {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if foundRoom == nil {
		return fmt.Errorf("room %s not found", roomID)
	}

	if foundRoom.Players == nil || len(foundRoom.Players) == 0 {
		s.log.Debug().Str("roomID", roomID).Msg("no players found for room")
		return nil
	}

	for _, player := range foundRoom.Players {
		if player == nil {
			continue
		}
		session := s.sessionManager.GetSession(player.TgID)
		if session != nil {
			session.UpdateState(newState)
		} else {
			s.log.Warn().Int64("tgID", player.TgID).Str("roomID", roomID).Msg("session not found for player")
		}
	}
	return nil
}

func (s *Service) updateAllActivePlayersState(roomID string, newState state.State) error {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if foundRoom == nil {
		return fmt.Errorf("room %s not found", roomID)
	}

	if foundRoom.Players == nil || len(foundRoom.Players) == 0 {
		s.log.Debug().Str("roomID", roomID).Msg("no players found for room")
		return nil
	}

	for _, player := range foundRoom.Round.ActivePlayers {
		if player == nil {
			continue
		}
		session := s.sessionManager.GetSession(player.TgID)
		if session != nil {
			session.UpdateState(newState)
		} else {
			s.log.Warn().Int64("tgID", player.TgID).Str("roomID", roomID).Msg("session not found for player")
		}
	}
	return nil
}

func (s *Service) updateAdminState(roomID string, state state.State) error {
	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		s.log.Error().Err(err).Str("roomID", roomID).Msg("failed to find room")
		return err
	}

	if foundRoom == nil {
		return fmt.Errorf("room %s not found", roomID)
	}

	if foundRoom.Admin == nil {
		return fmt.Errorf("no admin found for room %s", roomID)
	}

	session := s.sessionManager.GetSession(foundRoom.Admin.TgID)
	if session != nil {
		session.UpdateState(state)
	} else {
		s.log.Warn().Int64("tgID", foundRoom.Admin.TgID).Str("roomID", roomID).Msg("session not found for admin")
	}
	return nil
}
