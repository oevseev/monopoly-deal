package main

import (
	"github.com/oevseev/monopoly-deal/pkg/entities"
	"github.com/oevseev/monopoly-deal/pkg/game"
)

func main() {
	game.NewGame(entities.Settings{
		MovesPerTurn: 3,
		PlayerQueue: []entities.PlayerID{
			entities.PlayerID(1),
			entities.PlayerID(2),
			entities.PlayerID(3),
			entities.PlayerID(4),
		},
	}).Run()
}
