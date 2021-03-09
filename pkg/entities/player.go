package entities

type PlayerID int

type PlayerState struct {
	Hand     []Card
	Bank     []Card
	Property []Card
}

type PublicPlayerState struct {
	HandCount int
	Bank      []Card
	Property  []Card
}
