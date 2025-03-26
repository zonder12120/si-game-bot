package keyboard

import (
	"gopkg.in/telebot.v3"

	"github.com/zonder12120/tg-quiz/internal/telegram/render"
)

func Forbidden() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnReset,
	}, 1)
}

func MainMenu() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnJoinRoom,
		render.BtnCreateRoom,
	}, 2)
}

func NewRoundMenu() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnNewRound,
		render.BtnEndGame,
	}, 2)
}

func Back() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnBack,
	}, 1)
}

func StartRoundMenu() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnRoundStart,
		render.BtnBack,
	}, 2)
}

func EndRound() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnEndRound,
	}, 1)
}

func Leave() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnLeave,
	}, 1)
}

func Answer() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnAnswer,
	}, 1)
}

func ConfirmAnswer() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnYes,
		render.BtnNo,
	}, 2)
}

func ConfirmationMenu() *telebot.ReplyMarkup {
	return CreateKeyboard([]string{
		render.BtnConfirm,
		render.BtnCancel,
	}, 2)
}

func CreateKeyboard(buttons []string, columns int) *telebot.ReplyMarkup {
	kb := &telebot.ReplyMarkup{ResizeKeyboard: true}

	var rows []telebot.Row

	for i := 0; i < len(buttons); i += columns {
		var btns []telebot.Btn
		for j := i; j < i+columns && j < len(buttons); j++ {
			btns = append(btns, kb.Text(buttons[j]))
		}
		rows = append(rows, kb.Row(btns...))
	}

	kb.Reply(rows...)

	return kb
}
