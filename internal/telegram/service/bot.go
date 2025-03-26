package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/state"
)

type Bot struct {
	Bot             *telebot.Bot
	Log             *zerolog.Logger
	SessionsManager *state.Manager
}

type NewServiceParams struct {
	fx.In

	Bot             *telebot.Bot
	Log             *zerolog.Logger
	SessionsManager *state.Manager
}

func NewBot(params NewServiceParams) *Bot {
	return &Bot{
		Bot:             params.Bot,
		Log:             params.Log,
		SessionsManager: params.SessionsManager,
	}
}

func (b *Bot) SendMessage(telegramID int64, msg string, keyboard *telebot.ReplyMarkup) error {
	_, err := b.SendMsgAndGetInfo(telegramID, msg, keyboard)
	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) SendMsgAndGetInfo(telegramID int64, msg string, keyboard *telebot.ReplyMarkup) (*telebot.Message, error) {
	if b.Bot == nil {
		return nil, fmt.Errorf("bot instance is not initialized")
	}

	sendOptions := &telebot.SendOptions{
		ParseMode: telebot.ModeHTML,
	}

	if keyboard != nil {
		sendOptions.ReplyMarkup = keyboard
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var message *telebot.Message
	var sendErr error

	done := make(chan struct{})
	go func() {
		message, sendErr = b.Bot.Send(&telebot.User{ID: telegramID}, msg, sendOptions)
		close(done)
	}()

	select {
	case <-done:
		if sendErr != nil {
			b.Log.Error().
				Fields(struct{ TelegramID int64 }{telegramID}).
				Err(sendErr).
				Msg("failed to send message")
		}
		return message, sendErr
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (b *Bot) SendMsgAndGetInfoMdv2(telegramID int64, msg string, keyboard *telebot.ReplyMarkup) (*telebot.Message, error) {
	if b.Bot == nil {
		return nil, fmt.Errorf("bot instance is not initialized")
	}

	sendOptions := &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdownV2,
	}

	if keyboard != nil {
		sendOptions.ReplyMarkup = keyboard
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var message *telebot.Message
	var sendErr error

	done := make(chan struct{})
	go func() {
		message, sendErr = b.Bot.Send(&telebot.User{ID: telegramID}, msg, sendOptions)
		close(done)
	}()

	select {
	case <-done:
		if sendErr != nil {
			b.Log.Error().
				Fields(struct{ TelegramID int64 }{telegramID}).
				Err(sendErr).
				Msg("failed to send message")
		}
		return message, sendErr
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
