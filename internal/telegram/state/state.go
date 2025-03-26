package state

type State string

const (
	Idle State = "idle"

	// admin flow
	OnEnteredPoints  State = "on_entered_points"
	OnNewRoundMenu   State = "on_prepare_new_round"
	OnStartingRound  State = "on_starting_round"
	OnAdminRound     State = "on_admin_round"
	OnConfirmAnser   State = "on_confirm_anser"
	OnConfirmEndGame State = "on_confirm_stop"

	//player flow
	OnEnteredRoomID   State = "on_entered_room_id"
	OnWaitingNewRound State = "on_waiting_new_round"
	OnPlayingRound    State = "on_playing_round"
)
