package game

import (
	"container/ring"
	"errors"
	"fmt"
	"github.com/oevseev/monopoly-deal/pkg/entities"
)

type PlayerState struct {
	Hand     map[entities.CardID]struct{}
	Bank     map[entities.CardID]struct{}
	Property map[entities.CardID]struct{}
}

type Response struct {
	Transfers []entities.TransferProposal
	Confirmed bool
}

type Move struct {
	Cause     []entities.CardID
	Transfers map[entities.PlayerID][]entities.TransferProposal
	Responses map[entities.PlayerID]Response
}

func (m *Move) Empty() bool {
	return len(m.Cause) == 0
}

type State struct {
	cards        map[entities.CardID]entities.Card
	deck         []entities.CardID
	waste        []entities.CardID
	playerStates map[entities.PlayerID]PlayerState

	movesRemaining int
	move           Move

	settings    entities.Settings
	playerQueue *ring.Ring
}

func (s *State) idListToCardList(idList []entities.CardID) []entities.Card {
	result := make([]entities.Card, 0, len(idList))
	for _, id := range idList {
		result = append(result, s.cards[id])
	}
	return result
}

func (s *State) idSetToCardList(idSet map[entities.CardID]struct{}) []entities.Card {
	result := make([]entities.Card, 0, len(idSet))
	for cardID := range idSet {
		result = append(result, s.cards[cardID])
	}
	return result
}

func NewState(settings entities.Settings) *State {
	state := State{
		settings:       settings,
		movesRemaining: settings.MovesPerTurn,
		playerQueue:    ring.New(len(settings.PlayerQueue)),
	}
	for _, playerID := range settings.PlayerQueue {
		state.playerQueue.Value = playerID
		state.playerQueue = state.playerQueue.Next()
	}

	// TODO: initiate game state
	return &state
}

func (s *State) Public() entities.PublicGameState {
	var wasteTopCard entities.Card
	if len(s.waste) > 0 {
		wasteTopCard = s.cards[s.waste[0]]
	}

	playerStates := make(map[entities.PlayerID]entities.PublicPlayerState, len(s.playerStates))
	for playerID, playerState := range s.playerStates {
		playerStates[playerID] = entities.PublicPlayerState{
			HandCount: len(playerState.Hand),
			Bank:      s.idSetToCardList(playerState.Bank),
			Property:  s.idSetToCardList(playerState.Property),
		}
	}

	transfers := make(map[entities.PlayerID][]entities.Transfer, len(s.move.Transfers))
	for playerID, proposals := range s.move.Transfers {
		transfers[playerID] = make([]entities.Transfer, 0, len(proposals))
		for _, proposal := range proposals {
			transfers[playerID] = append(transfers[playerID], entities.Transfer{
				Card:   s.cards[proposal.CardID],
				Target: proposal.Target,
			})
		}
	}

	responses := make(map[entities.PlayerID]entities.Response, len(s.move.Responses))
	for playerID, internalResponse := range s.move.Responses {
		transfers := make([]entities.Transfer, 0, len(internalResponse.Transfers))
		for _, transfer := range internalResponse.Transfers {
			transfers = append(transfers, entities.Transfer{
				Card:   s.cards[transfer.CardID],
				Target: transfer.Target,
			})
		}
		responses[playerID] = entities.Response{
			Transfers: transfers,
			Confirmed: internalResponse.Confirmed,
		}
	}

	move := entities.Move{
		Cause:     s.idListToCardList(s.move.Cause),
		Transfers: transfers,
		Responses: responses,
	}

	state := entities.PublicGameState{
		DeckCardCount:      len(s.deck),
		WasteTopCard:       wasteTopCard,
		PublicPlayerStates: playerStates,
		CurrentPlayer:      s.playerQueue.Value.(entities.PlayerID),
		MovesRemaining:     s.movesRemaining,
		Move:               move,
	}
	return state
}

func (s *State) NextTurn() {
	s.playerQueue = s.playerQueue.Next()
	s.movesRemaining = s.settings.MovesPerTurn
}

func (s *State) CurrentPlayerID() entities.PlayerID {
	return s.playerQueue.Value.(entities.PlayerID)
}

func (s *State) PlayerState(id entities.PlayerID) entities.PlayerState {
	playerState := s.playerStates[id]
	return entities.PlayerState{
		Hand:     s.idSetToCardList(playerState.Hand),
		Bank:     s.idSetToCardList(playerState.Bank),
		Property: s.idSetToCardList(playerState.Property),
	}
}

func (s *State) commenceMove() {
	// TODO: implement commenceMove
}

func (s *State) postResponse() {
	allConfirmed := true
	for playerID := range s.move.Transfers {
		response, ok := s.move.Responses[playerID]
		if !ok || !response.Confirmed {
			allConfirmed = false
			break
		}
	}
	if allConfirmed {
		s.commenceMove()
		if s.movesRemaining == 0 {
			s.NextTurn()
		}
	}
}

func (s *State) ProposeMove(p entities.MoveProposal) error {
	if !s.move.Empty() {
		return errors.New("move already in progress")
	}
	if len(p.Cause) == 0 {
		return errors.New("empty cause")
	}
	s.move = Move{
		Cause:     p.Cause,
		Transfers: p.Transfers,
		Responses: make(map[entities.PlayerID]Response),
	}
	s.movesRemaining--
	return nil
}

func (s *State) Respond(playerID entities.PlayerID, transfers []entities.TransferProposal) error {
	if s.move.Empty() {
		return errors.New("no move in progress")
	}
	if _, ok := s.move.Responses[playerID]; ok {
		return errors.New("already responded")
	}
	s.move.Responses[playerID] = Response{
		Transfers: transfers,
		Confirmed: true,
	}
	s.postResponse()
	return nil
}

func (s *State) SayNo(playerID entities.PlayerID) error {
	if s.move.Empty() {
		return errors.New("no move in progress")
	}
	if _, ok := s.move.Responses[playerID]; ok {
		return fmt.Errorf("player %d already responded", playerID)
	}
	s.move.Responses[playerID] = Response{
		Transfers: nil,
		Confirmed: false,
	}
	return nil
}

func (s *State) AcceptNo(playerID entities.PlayerID) error {
	if s.move.Empty() {
		return errors.New("no move in progress")
	}
	response, ok := s.move.Responses[playerID]
	if !ok {
		return fmt.Errorf("no response from player %d yet", playerID)
	}
	if response.Confirmed {
		return fmt.Errorf("response for player %d already confirmed", playerID)
	}
	s.move.Responses[playerID] = Response{
		Transfers: response.Transfers,
		Confirmed: true,
	}
	s.postResponse()
	return nil
}

func (s *State) DeclineNo(playerID entities.PlayerID) error {
	if s.move.Empty() {
		return errors.New("no move in progress")
	}
	response, ok := s.move.Responses[playerID]
	if !ok {
		return fmt.Errorf("no response from player %d yet", playerID)
	}
	if response.Confirmed {
		return fmt.Errorf("response for player %d already confirmed", playerID)
	}
	delete(s.move.Responses, playerID)
	return nil
}

func (s *State) Draw(count int) {
	// TODO: implement Draw
}

func (s *State) Discard(id entities.CardID) {
	// TODO: implement Discard
}
