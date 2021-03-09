package entities

type CardID int
type CardKindID int
type CardLocationType string

const (
	Deck           CardLocationType = "Deck"
	Waste                           = "Waste"
	PlayerHand                      = "PlayerHand"
	PlayerBank                      = "PlayerBank"
	PlayerProperty                  = "PlayerProperty"
)

type CardLocation struct {
	Type     CardLocationType
	PlayerID PlayerID
}

type Card struct {
	ID       CardID
	Location CardLocation

	// TODO: add card kind and metadata
}
