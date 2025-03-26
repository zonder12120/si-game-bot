package access

import "github.com/zonder12120/tg-quiz/internal/telegram/user"

type Command string

const (
	CmdNewRound      Command = "new_round"
	CmdEndGame       Command = "end_game"
	CmdStartRound    Command = "start_round"
	CmdConfirmAnswer Command = "confirm_answer"
)

var accessMap = map[Command][]user.Role{
	CmdNewRound:      {user.Admin},
	CmdEndGame:       {user.Admin},
	CmdStartRound:    {user.Admin},
	CmdConfirmAnswer: {user.Admin},
}

func HasAccess(group user.Role, command Command) bool {
	for _, allowedGroup := range accessMap[command] {
		if allowedGroup == group {
			return true
		}
	}
	return false
}
