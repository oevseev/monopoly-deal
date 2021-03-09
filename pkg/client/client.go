package client

import "github.com/oevseev/monopoly-deal/pkg/entities"

type Client interface {
	RequestMoveProposal(entities.PublicGameState, entities.PlayerState, chan<- entities.MoveProposal)
	RequestResponse(entities.PublicGameState, entities.PlayerState, chan<- entities.ResponseProposal)
	RequestSayNoConfirmation(entities.PublicGameState, entities.PlayerState, chan<- bool)
}
