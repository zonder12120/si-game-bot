package game

import (
	"context"
	"fmt"
	"time"

	"github.com/zonder12120/tg-quiz/internal/telegram/render"
)

type TimerType string

const (
	RoundTimer  TimerType = "round"
	AnswerTimer TimerType = "answer"
)

type Timer struct {
	ctx        context.Context
	cancel     context.CancelFunc
	duration   time.Duration
	remaining  time.Duration
	paused     bool
	messageIDs map[int64]int
}

func (s *Service) StartTimer(roomID string, timerType TimerType, duration time.Duration, tgIDs []int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.timers[roomID]; !exists {
		s.timers[roomID] = make(map[TimerType]*Timer)
	}

	if existing, ok := s.timers[roomID][timerType]; ok {
		existing.cancel()
	}

	ctx, cancel := context.WithCancel(context.Background())
	messageText := fmt.Sprintf(render.MsgTimer, int(duration.Seconds()))
	messageIDs, _ := s.notifyService.NotifyUsers(tgIDs, messageText, nil)

	s.timers[roomID][timerType] = &Timer{
		ctx:        ctx,
		cancel:     cancel,
		duration:   duration,
		remaining:  duration,
		paused:     false,
		messageIDs: messageIDs,
	}

	go s.runTimer(roomID, timerType, ctx, messageIDs, tgIDs)

	s.log.Debug().
		Str("roomID", roomID).
		Str("timerType", string(timerType)).
		Msg("Start timer")
}

func (s *Service) runTimer(roomID string, timerType TimerType, ctx context.Context, messageIDs map[int64]int, tgIDs []int64) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			timer, exists := s.timers[roomID][timerType]
			if !exists || timer.paused {
				s.mu.Unlock()
				return
			}
			timer.remaining -= 1 * time.Second
			if timer.remaining > 0 {
				text := fmt.Sprintf(render.MsgTimer, int(timer.remaining.Seconds()))
				s.notifyService.UpdateMessages(messageIDs, text)
				s.mu.Unlock()
				continue
			}
			delete(s.timers[roomID], timerType)
			if len(s.timers[roomID]) == 0 {
				delete(s.timers, roomID)
			}
			s.mu.Unlock()

			_, err := s.notifyService.NotifyUsers(tgIDs, render.MsgEndTime, nil)
			if err != nil {
				return
			}

			if timerType == AnswerTimer {
				err = s.incorrectAnswer(roomID)
				if err != nil {
					return
				}
			}

			_ = s.EndRound(roomID)
			return

		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) PauseTimer(roomID string, timerType TimerType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if timers, exists := s.timers[roomID]; exists {
		if timer, ok := timers[timerType]; ok {
			timer.cancel()
			timer.paused = true

			s.log.Debug().
				Str("roomID", roomID).
				Str("timerType", string(timerType)).
				Msg("Paused timer")
		}
	}
}

func (s *Service) ResumeTimer(roomID string, timerType TimerType, tgIDs []int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if timers, exists := s.timers[roomID]; exists {
		if timer, ok := timers[timerType]; ok && timer.paused {
			ctx, cancel := context.WithCancel(context.Background())
			timer.ctx = ctx
			timer.cancel = cancel
			timer.paused = false

			messageText := fmt.Sprintf(render.MsgTimer, int(timer.remaining.Seconds()))
			messageIDs, _ := s.notifyService.NotifyUsers(tgIDs, messageText, nil)
			timer.messageIDs = messageIDs

			go s.runTimer(roomID, timerType, ctx, messageIDs, tgIDs)

			s.log.Debug().
				Str("roomID", roomID).
				Str("timerType", string(timerType)).
				Msg("Resume timer")
		}
	}
}

func (s *Service) StopTimer(roomID string, timerType TimerType) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if timers, exists := s.timers[roomID]; exists {
		if timer, ok := timers[timerType]; ok {
			timer.cancel()

			delete(timers, timerType)

			if len(timers) == 0 {
				delete(s.timers, roomID)
			}

			s.log.Debug().
				Str("roomID", roomID).
				Str("timerType", string(timerType)).
				Msg("Timer stopped")
		}
	}
}

func (s *Service) StopAllTimers(roomID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if timers, exists := s.timers[roomID]; exists {
		for _, timer := range timers {
			timer.cancel()
		}
		delete(s.timers, roomID)
		s.log.Debug().
			Str("roomID", roomID).
			Msg("All Timer stopped")
	}
}
