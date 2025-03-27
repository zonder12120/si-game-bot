package render

const (

	// cancel
	MsgCancel = "❗ Действие отменено."

	// start
	MsgStartCommand = "Выбери действие:"

	// unknown command
	MsgUnknownCommand = "❌ Неизвестная команда или неверное состояние."

	// access middleware
	ErrMessageTextForbidden = "⛔ У вас нет доступа к этой команде."

	// errors
	ErrUnexpected = "❌ Неизвестная ошибка: %s"

	// confirmation
	MsgAreYouSure = "Вы уверены?"

	// admin flow
	MsgRoomCreated   = "✅ Комната создана, её id: ```%s```\nСкопируйте id и отправьте друзьям\nНажав кнопку \"Новый раунд\", вы запускаете игру, и к ней нельзя будет подключиться"
	MsgCreateRound   = "Создайте новый раунд или завершите игру!"
	MsgEnterPoints   = "Введите количество очков за вопрос: "
	MsgRoundCreated  = "✅ Раунд создан, можете начинать:"
	MsgRoundStarted  = "✅ Раунд начался!"
	MsgConfirmAnswer = "Оцените ответ игрока"

	// player flow
	MsgEnterRoomID     = "Введите ID комнаты (его можно получить у того, кто создавал комнату):"
	MsgJoinedInTheRoom = "✅ Вы успешно присоеденились к игре!\nОжидание начала..."
	MsgLeaveGame       = "⚠️ Вы вышли из игры"

	// notifications
	MsgPlayerJoined          = "%s присоединился(лась) к игре 👋"
	MsgWaitAnswering         = "Ожидайте этап ответа игроков"
	MsgNewRoundStarted       = "📢 Начался новый раунд, слушайте вопрос!"
	MsgCanAnswer             = "🟢 Можете отвечать!"
	MsgWaitAnswer            = "Отвечает игрок %s"
	MsgCorrectAnswer         = "✅ Верный ответ!"
	MsgIncorrectAnswer       = "❌ Не верный ответ!"
	MsgEndRound              = "🔚 Раунд завершён! Результаты:"
	MsgEndGame               = "🥳 Игра завершена! Результаты:"
	MsgTimer                 = "Осталось %d %s"
	MsgEndTime               = "🔴 Время вышло!"
	MsgCorrectAnswerResult   = "+%d %s"
	MsgIncorrectAnswerResult = "%d %s"
	MsgPlayerLeft            = "%s покинул игру"
	MsgCantJoin              = "Невозможно подключиться к игре"

	MsgNoPlayers = "\nНет игроков!"
	MsgResult    = "\nИгрок %s: %d %s"
	MsgWinner    = "\n🏆 Победитель: "
	MsgWinners   = "\n🏆 Победители: "
)
