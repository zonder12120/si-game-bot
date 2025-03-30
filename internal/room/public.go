package room

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/zonder12120/tg-quiz/internal/telegram/render"
	"github.com/zonder12120/tg-quiz/internal/telegram/user"
)

func (m *Manager) GetRoom(roomID string) (*RoomSnapshot, error) {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return nil, fmt.Errorf("комната не найдена")
	}
	room := value.(*Room)

	room.mu.RLock()
	defer room.mu.RUnlock()

	return createSnapshot(room), nil
}

func (m *Manager) GetMember(roomID string, tgID int64) (*Member, error) {
	roomSnapshot, err := m.GetRoom(roomID)
	if err != nil {
		return nil, err
	}

	foundPlayer, playerOk := roomSnapshot.Players[tgID]
	adminOk := roomSnapshot.Admin.TgID == tgID

	if adminOk {
		return roomSnapshot.Admin, nil
	}
	if playerOk {
		return foundPlayer, nil
	}
	return nil, fmt.Errorf("не найден указанный игрок")

}

func (m *Manager) CreateRoom(adminTgID int64) string {
	id := uuid.New().String()
	room := &Room{
		Admin: &Member{
			TgID: adminTgID,
			Role: user.Admin,
		},
		Players: make(map[int64]*Member),
		Status:  waiting,
	}
	room.updateTimestamp()
	m.rooms.Store(id, room)
	return id
}

func (m *Manager) JoinRoom(roomID string, tgID int64, name string) error {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("комната не найдена")
	}
	foundRoom := value.(*Room)

	foundRoom.mu.Lock()
	defer foundRoom.mu.Unlock()

	if _, ok = m.rooms.Load(roomID); !ok {
		return fmt.Errorf("комната удалена")
	}

	if foundRoom.Status != waiting {
		return fmt.Errorf("игра уже идёт")
	}

	foundRoom.Players[tgID] = &Member{
		TgID: tgID,
		Name: name,
		Role: user.Player,
	}
	foundRoom.updateTimestamp()
	return nil
}

func (m *Manager) EndGame(roomID string) error {
	m.rooms.Delete(roomID)
	return nil
}

func (r *Room) updateTimestamp() {
	r.updatedAt.Store(time.Now())
}

func (m *Manager) LeaveRoom(roomID string, tgID int64) error {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("комната не найдена")
	}
	foundRoom := value.(*Room)

	foundRoom.mu.Lock()
	defer foundRoom.mu.Unlock()

	foundPlayer, ok := foundRoom.Players[tgID]
	if !ok {
		return fmt.Errorf("не найден указанный игрок")
	}

	delete(foundRoom.Players, foundPlayer.TgID)
	foundRoom.updateTimestamp()

	return nil
}

func (m *Manager) StartGame(roomID string) error {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("комната не найдена")
	}
	foundRoom := value.(*Room)

	foundRoom.mu.Lock()
	defer foundRoom.mu.Unlock()

	if foundRoom.Status == playing {
		return nil
	}

	foundRoom.Status = playing
	foundRoom.updateTimestamp()
	return nil
}

func (m *Manager) NewRound(roomID string, points int) error {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("комната не найдена")
	}
	foundRoom := value.(*Room)

	foundRoom.mu.Lock()
	defer foundRoom.mu.Unlock()

	if foundRoom.Status != playing {
		foundRoom.mu.Unlock()
		err := m.StartGame(roomID)
		foundRoom.mu.Lock()
		if err != nil {
			return err
		}
	}

	activePlayers := make(map[int64]*Member, len(foundRoom.Players))
	for id, p := range foundRoom.Players {
		activePlayers[id] = &Member{
			TgID:   p.TgID,
			Name:   p.Name,
			Role:   p.Role,
			Points: p.Points,
		}
	}

	foundRoom.Round = &Round{
		ActivePlayers:   activePlayers,
		Points:          points,
		AnsweringPlayer: 0,
	}
	foundRoom.updateTimestamp()
	return nil
}

func (m *Manager) StartAnswer(roomID string, tgID int64) error {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("комната не найдена")
	}
	foundRoom := value.(*Room)

	foundRoom.mu.Lock()
	defer foundRoom.mu.Unlock()

	if foundRoom.Status != playing {
		return fmt.Errorf("нельзя давать ответ, пока не началась игра")
	}

	if foundRoom.Round == nil {
		return fmt.Errorf("нельзя давать ответ, пока не начался раунд")
	}

	foundPlayer, ok := foundRoom.Players[tgID]
	if !ok {
		return fmt.Errorf("не найден указанный игрок")
	}

	if foundRoom.Round.AnsweringPlayer != 0 {
		return nil
	}

	foundRoom.updateTimestamp()
	foundRoom.Round.AnsweringPlayer = foundPlayer.TgID

	return nil
}

func (m *Manager) ResultAnswer(roomID string, isRight bool) error {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return fmt.Errorf("комната не найдена")
	}
	foundRoom := value.(*Room)

	foundRoom.mu.Lock()
	defer foundRoom.mu.Unlock()

	if foundRoom.Status != playing {
		return fmt.Errorf("нельзя давать ответ, пока не началась игра")
	}

	if foundRoom.Round == nil {
		return fmt.Errorf("нельзя давать ответ, пока не начался раунд")
	}

	if foundRoom.Round.AnsweringPlayer == 0 {
		return fmt.Errorf("некому засчитывать ответ")
	}

	points := foundRoom.Round.Points
	answeringPlayer := foundRoom.Round.AnsweringPlayer

	if !isRight {
		points *= -1
		delete(foundRoom.Round.ActivePlayers, answeringPlayer)
	}

	foundRoom.Players[answeringPlayer].Points += points
	foundRoom.Round.AnsweringPlayer = 0
	foundRoom.updateTimestamp()

	return nil
}

func (m *Manager) GetTopPlayersNames(roomID string) ([]string, error) {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return nil, fmt.Errorf("комната не найдена")
	}
	foundRoom := value.(*Room)

	foundRoom.mu.RLock()
	defer foundRoom.mu.RUnlock()

	if len(foundRoom.Players) == 0 {
		return []string{
			render.MsgNoPlayers,
		}, nil
	}

	var maxPoints int
	topPlayers := make([]string, 0)

	first := true
	for _, p := range foundRoom.Players {
		if first {
			maxPoints = p.Points
			topPlayers = append(topPlayers, p.Name)
			first = false
			continue
		}

		switch {
		case p.Points > maxPoints:
			maxPoints = p.Points
			topPlayers = []string{p.Name}
		case p.Points == maxPoints:
			topPlayers = append(topPlayers, p.Name)
		}
	}

	return topPlayers, nil
}
