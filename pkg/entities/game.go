package entities

type PublicGameState struct {
	DeckCardCount      int
	WasteTopCard       Card
	PublicPlayerStates map[PlayerID]PublicPlayerState
	CurrentPlayer      PlayerID
	MovesRemaining     int
	Move               Move
}
