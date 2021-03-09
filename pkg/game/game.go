package game

import (
	"github.com/oevseev/monopoly-deal/internal/game"
	"github.com/oevseev/monopoly-deal/pkg/entities"
)

type Game interface {
	Run()
}

func NewGame(settings entities.Settings) Game {
	return &game.Game{
		Settings: settings,
		State:    game.NewState(settings),
	}
}
