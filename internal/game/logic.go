package game

import (
	"github.com/oevseev/monopoly-deal/pkg/client"
	"github.com/oevseev/monopoly-deal/pkg/entities"
)

type Game struct {
	Settings entities.Settings
	State    *State
	Clients  map[entities.PlayerID]client.Client
}

func (g *Game) Run() {
	// TODO: implement Run
}

func (g *Game) moveIteration() {
	playerID := g.State.CurrentPlayerID()

	proposalCh := make(chan entities.MoveProposal, 1)
	go g.Clients[playerID].RequestMoveProposal(
		g.State.Public(),
		g.State.PlayerState(playerID),
		proposalCh)

	// TODO: finish the implementation of moveIteration
}
