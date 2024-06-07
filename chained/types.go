package chained

const NumNodes = 3
const Generic MessageType = "Generic"

type MessageType string
type ViewNumber = int
type VotePower = int
type BlockHash = string

type Proof struct {
	BlockHash BlockHash
	Proof     VotePower
}

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

type Node struct {
	id              int
	myViewNumber    ViewNumber
	hashes          map[BlockHash]VotePower
	voteBlockHash   BlockHash
	futureBlockHash BlockHash
	blockchain      *Blockchain
}
