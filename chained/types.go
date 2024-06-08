package chained

import "strconv"

const NumNodes = 2
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
	Send(int, Msg)
}

type Node struct {
	id              int
	myViewNumber    ViewNumber
	hashes          map[BlockHash]VotePower
	voteBlockHash   BlockHash
	futureBlockHash BlockHash
	blockchain      *Blockchain
}

type Block struct {
	id int
}

func newBlock(id int) *Block {
	return &Block{id: id}
}

func (b *Block) Hash() string {
	return strconv.Itoa(b.id)
}
