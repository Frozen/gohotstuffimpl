package chained

type MessageType string
type ViewNumber = int
type VotePower = int
type BlockHash = string
type Proof struct {
	BlockHash BlockHash
	Proof     []byte
}

const Generic MessageType = "Generic"

type Msg struct {
	Type          MessageType
	ViewNumber    ViewNumber
	SenderID      int // node or public key
	VoteBlockHash BlockHash
	Block         bool
	Proof         *Proof
}

type Responder interface {
	Broadcast(Msg)
	ToLeader(Msg)
}
