package room

import (
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/zonder12120/tg-quiz/internal/config"
)

func NewManager(config config.RoomManager, logger *zerolog.Logger) *Manager {
	m := &Manager{
		ttl:    config.TTL,
		stopCh: make(chan struct{}),
		wg:     sync.WaitGroup{},
		log:    logger,
	}
	go m.cleanupWorker()
	return m
}

func (m *Manager) cleanupWorker() {
	m.wg.Add(1)
	defer m.wg.Done()

	ticker := time.NewTicker(m.ttl)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.log.Debug().Msg("Cleaning up!")
			var toDelete []string
			m.rooms.Range(func(key, value interface{}) bool {
				room := value.(*Room)
				ts := room.updatedAt.Load().(time.Time)
				if time.Since(ts) > m.ttl {
					toDelete = append(toDelete, key.(string))
				}
				return true
			})
			for _, id := range toDelete {
				m.rooms.Delete(id)
				m.log.Info().Str("roomID", id).Msg("removing inactive room")
			}
		case <-m.stopCh:
			return
		}
	}
}

func (m *Manager) Stop() {
	m.log.Info().Msg("Stopping room manager")
	close(m.stopCh)
	m.wg.Wait()
	m.log.Info().Msg("Room manager stopped")
}

func (m *Manager) getRoom(roomID string) (*RoomSnapshot, error) {
	value, ok := m.rooms.Load(roomID)
	if !ok {
		return nil, fmt.Errorf("комната не найдена")
	}
	room := value.(*Room)

	room.mu.RLock()
	defer room.mu.RUnlock()

	return createSnapshot(room), nil
}

func createSnapshot(room *Room) *RoomSnapshot {
	snapshot := &RoomSnapshot{
		Admin:   room.Admin,
		Players: make(map[int64]*Member, len(room.Players)),
		Status:  room.Status,
	}

	for id, p := range room.Players {
		snapshot.Players[id] = &Member{
			TgID:   p.TgID,
			Name:   p.Name,
			Role:   p.Role,
			Points: p.Points,
		}
	}

	if room.Round != nil {
		snapshot.Round = &RoundSnapshot{
			ActivePlayers:   make(map[int64]*Member, len(room.Round.ActivePlayers)),
			AnsweringPlayer: room.Round.AnsweringPlayer,
			Points:          room.Round.Points,
		}
		for id, p := range room.Round.ActivePlayers {
			snapshot.Round.ActivePlayers[id] = &Member{
				TgID:   p.TgID,
				Name:   p.Name,
				Role:   p.Role,
				Points: p.Points,
			}
		}
	}

	return snapshot
}
