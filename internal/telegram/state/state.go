package state

type State string

const (
	Idle State = "idle"

	// admin flow
	OnEnteredPoints  State = "on_entered_points"
	OnNewRoundMenu   State = "on_prepare_new_round"
	OnStartingRound  State = "on_starting_round"
	OnAdminRound     State = "on_admin_round"
	OnConfirmAnswer  State = "on_confirm_answer"
	OnConfirmEndGame State = "on_confirm_end_game"

	//player flow
	OnEnteredRoomID   State = "on_entered_room_id"
	OnWaitingNewRound State = "on_waiting_new_round"
	OnPlayingRound    State = "on_playing_round"
	OnLeavingRoom     State = "on_leaving_round"
)
