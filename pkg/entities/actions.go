package entities

type TransferProposal struct {
	CardID CardID
	Target CardLocation
}

type Transfer struct {
	Card   Card
	Target CardLocation
}

type Response struct {
	Transfers []Transfer
	Confirmed bool
}

type ResponseProposal struct {
	SayNo     bool
	Transfers []TransferProposal
}

type Move struct {
	Cause     []Card
	Transfers map[PlayerID][]Transfer
	Responses map[PlayerID]Response
}

type MoveProposal struct {
	Cause     []CardID
	Transfers map[PlayerID][]TransferProposal
}

func (m *Move) Empty() bool {
	return len(m.Cause) == 0
}
