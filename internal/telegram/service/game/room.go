package game

import (
	"fmt"
	"strings"

	"github.com/zonder12120/tg-quiz/internal/telegram/keyboard"
	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

func (s *Service) CreateRoom(tgID int64) string {
	return s.roomManager.CreateRoom(tgID)
}

func (s *Service) JoinRoom(roomID string, tgID int64, fullName string) error {
	err := s.roomManager.JoinRoom(roomID, tgID, fullName)
	if err != nil {
		return err
	}
	_, err = s.notifyAllMembers(roomID, fmt.Sprintf(render.MsgPlayerJoined, fullName), nil)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) LeaveRoom(roomID string, tgID int64) error {
	player, err := s.GetPlayer(roomID, tgID)
	if err != nil {
		return err
	}

	err = s.roomManager.LeaveRoom(roomID, tgID)
	if err != nil {
		return err
	}
	tgIDs, err := s.getAllMembersIDs(roomID)
	if err != nil {
		return err
	}

	_, err = s.notifyService.NotifyUsers(
		tgIDs,
		fmt.Sprintf(render.MsgPlayerLeft, player.Name),
		nil,
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) NewRound(roomID string, points int) error {
	if points < 0 {
		points *= -1
	}
	err := s.roomManager.NewRound(roomID, points)
	if err != nil {
		return err
	}
	_, err = s.notifyAllPlayers(roomID, render.MsgNewRoundStarted, keyboard.Leave())
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) StartRound(roomID string) error {
	if err := s.updateAllPlayersState(roomID, state.OnPlayingRound); err != nil {
		s.log.Error().Err(err).Str("roomID", roomID).Msg("failed to update player states")
	}

	if _, err := s.notifyAllPlayers(roomID, render.MsgCanAnswer, keyboard.Answer()); err != nil {
		s.log.Error().Err(err).Str("roomID", roomID).Msg("failed to notify players")
	}

	tgIDs, err := s.getAllMembersIDs(roomID)
	if err != nil {
		return err
	}

	s.StartTimer(roomID, RoundTimer, s.gameConfig.RoundTTL, tgIDs)

	return nil
}

func (s *Service) StartAnswer(roomID string, tgID int64) error {
	err := s.roomManager.StartAnswer(roomID, tgID)
	if err != nil {
		return err
	}

	s.PauseTimer(roomID, RoundTimer)

	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return err
	}

	playerName := foundRoom.Players[tgID].Name

	_, err = s.notifyAllMembers(roomID, fmt.Sprintf(render.MsgWaitAnswer, playerName), keyboard.Leave())
	if err != nil {
		return err
	}

	err = s.updateAllPlayersState(roomID, state.OnWaitingNewRound)
	if err != nil {
		return err
	}

	err = s.updateAdminState(roomID, state.OnConfirmAnswer)
	if err != nil {
		return err
	}

	err = s.notifyAdmin(
		roomID,
		render.MsgConfirmAnswer,
		keyboard.ConfirmAnswer(),
	)
	if err != nil {
		return err
	}

	tgIDs, err := s.getAllMembersIDs(roomID)
	if err != nil {
		return err
	}

	s.StartTimer(roomID, AnswerTimer, s.gameConfig.AnswerTTL, tgIDs)

	return nil
}

func (s *Service) ResultAnswer(roomID string, isRight bool) error {

	if isRight {
		err := s.correctAnswer(roomID)
		if err != nil {
			return err
		}
	} else {
		err := s.incorrectAnswer(roomID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) correctAnswer(roomID string) error {
	answeringPlayer, err := s.getAnsweringPlayer(roomID)
	if err != nil {
		return err
	}

	err = s.roomManager.ResultAnswer(roomID, true)
	if err != nil {
		return err
	}

	points, err := s.getRoundPoints(roomID)
	if err != nil {
		return err
	}

	_, errNotify := s.notifyAllMembers(roomID, render.MsgCorrectAnswer, nil)
	if errNotify != nil {
		return errNotify
	}
	err = s.botService.SendMessage(
		answeringPlayer,
		fmt.Sprintf(render.MsgCorrectAnswerResult, points, getPointsWord(points)),
		nil,
	)
	if err != nil {
		return err
	}

	err = s.EndRound(roomID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) incorrectAnswer(roomID string) error {
	s.StopTimer(roomID, AnswerTimer)

	answeringPlayer, err := s.getAnsweringPlayer(roomID)
	if err != nil {
		return err
	}

	err = s.roomManager.ResultAnswer(roomID, false)
	if err != nil {
		return err
	}

	points, err := s.getRoundPoints(roomID)
	if err != nil {
		return err
	}

	points *= -1

	_, errNotify := s.notifyAllPlayers(roomID, render.MsgIncorrectAnswer, nil)
	if errNotify != nil {
		return errNotify
	}

	err = s.botService.SendMessage(
		answeringPlayer,
		fmt.Sprintf(render.MsgIncorrectAnswerResult, points, getPointsWord(points)),
		nil,
	)
	if err != nil {
		return err
	}

	err = s.updateAllActivePlayersState(roomID, state.OnPlayingRound)
	if err != nil {
		return err
	}

	err = s.notifyAllActivePlayers(roomID, render.MsgCanAnswer, keyboard.Answer())
	if err != nil {
		return err
	}

	err = s.updateAdminState(roomID, state.OnAdminRound)
	if err != nil {
		return err
	}

	errNotify = s.notifyAdmin(roomID, render.MsgIncorrectAnswer, keyboard.EndRound())
	if errNotify != nil {
		return errNotify
	}

	tgIDs, errGetIDs := s.getAllMembersIDs(roomID)
	if errGetIDs != nil {
		return errGetIDs
	}

	s.ResumeTimer(roomID, RoundTimer, tgIDs)

	return nil
}

func (s *Service) EndRound(roomID string) error {
	s.StopAllTimers(roomID)

	results, err := s.getGameResults(roomID)
	if err != nil {
		return err
	}

	msg := render.MsgEndRound + results

	err = s.updateAllPlayersState(roomID, state.OnWaitingNewRound)
	if err != nil {
		return err
	}

	_, err = s.notifyAllPlayers(roomID, msg, keyboard.Leave())
	if err != nil {
		return err
	}

	err = s.updateAdminState(roomID, state.OnNewRoundMenu)
	if err != nil {
		return err
	}
	err = s.notifyAdmin(roomID, msg, keyboard.NewRoundMenu())
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) EndGame(roomID string) error {
	results, err := s.getGameResults(roomID)
	if err != nil {
		return err
	}

	msg := render.MsgEndGame + results

	topPlayers, err := s.roomManager.GetTopPlayersNames(roomID)
	if err != nil {
		return err
	}
	switch len(topPlayers) {
	case 0:
		s.log.Error().Err(err).Str("roomID", roomID).Msg("no top players found")
		return fmt.Errorf("не найдено победителей")
	case 1:
		msg += render.MsgWinner + topPlayers[0]
	default:
		msg += render.MsgWinners
		msg += strings.Join(topPlayers, ", ")
	}

	_, err = s.notifyAllMembers(roomID, msg, keyboard.MainMenu())
	if err != nil {
		s.log.Error().Err(err).Str("roomID", roomID).Msg("failed to notify all players")
		return err
	}

	err = s.updateAllPlayersState(roomID, state.Idle)
	if err != nil {
		s.log.Error().Err(err).Str("roomID", roomID).Msg("failed to update all players")
	}

	err = s.updateAdminState(roomID, state.Idle)
	if err != nil {
		s.log.Error().Err(err).Str("roomID", roomID).Msg("failed to update admin state")
	}

	err = s.roomManager.EndGame(roomID)
	if err != nil {
		s.log.Error().Err(err).Str("roomID", roomID).Msg("failed to stop game")
	}

	return err
}

func (s *Service) getGameResults(roomID string) (string, error) {
	var builder strings.Builder

	foundRoom, err := s.roomManager.GetRoom(roomID)
	if err != nil {
		return "", err
	}

	if len(foundRoom.Players) == 0 {
		return render.MsgNoPlayers, nil
	}

	for _, player := range foundRoom.Players {
		builder.WriteString(fmt.Sprintf(render.MsgResult, player.Name, player.Points, getPointsWord(player.Points)))
	}

	return builder.String(), nil
}

func getPointsWord(points int) string {
	n := abs(points)
	lastDigit := n % 10
	lastTwoDigits := n % 100

	switch {
	case lastDigit == 1 && lastTwoDigits != 11:
		return "очко"
	case lastDigit >= 2 && lastDigit <= 4 && (lastTwoDigits < 10 || lastTwoDigits >= 20):
		return "очка"
	default:
		return "очков"
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (s *Service) getAnsweringPlayer(roomID string) (int64, error) {
	foundRoom, foundErr := s.roomManager.GetRoom(roomID)
	if foundErr != nil {
		return 0, foundErr
	}

	answeringPlayer := foundRoom.Round.AnsweringPlayer
	if answeringPlayer == 0 {
		return 0, fmt.Errorf("в комнате нет отвечающих игроков")
	}

	return foundRoom.Round.AnsweringPlayer, nil
}

func (s *Service) getRoundPoints(roomID string) (int, error) {
	foundRoom, foundErr := s.roomManager.GetRoom(roomID)
	if foundErr != nil {
		return 0, foundErr
	}

	return foundRoom.Round.Points, nil
}
