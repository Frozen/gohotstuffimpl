package chained

import (
	"fmt"
	"strconv"
)

func NewNode(id int) *Node {
	return &Node{
		id:            id,
		myViewNumber:  0,
		blockchain:    NewBlockchain(),
		voteBlockHash: strconv.Itoa(0),
		hashes: map[BlockHash]VotePower{
			strconv.Itoa(0): 1,
		},
	}
}

func (n *Node) Genesis() Msg {
	return Msg{
		Type:          Generic,
		ViewNumber:    0,
		SenderID:      n.id,
		VoteBlockHash: strconv.Itoa(0),
		Block:         true,
		Proof:         nil,
	}
}

func (n *Node) Apply(msg Msg, responder Responder) {
	switch { // I should be leader
	case n.myViewNumber == msg.ViewNumber: // here we collect sigs
		if msg.VoteBlockHash != n.voteBlockHash {
			panic(fmt.Sprintf("expected vote block hash '%s', got '%s'", n.voteBlockHash, msg.VoteBlockHash))
		}

		// TODO: check time as well

		if n.vote(msg.VoteBlockHash) == NumNodes { // 100%
			n.blockchain.AddBlock(msg.VoteBlockHash)
			newBlock := proposeNewBlock(n.myViewNumber)
			n.futureBlockHash = newBlock.Hash()

			responder.Broadcast(Msg{
				Type:          Generic,
				ViewNumber:    msg.ViewNumber + 1,
				SenderID:      n.id,
				Block:         isLeader(msg.ViewNumber+1, n.id),
				VoteBlockHash: n.futureBlockHash,
				Proof: &Proof{
					BlockHash: n.voteBlockHash,
					Proof:     n.hashes[n.voteBlockHash],
				},
			})

			n.myViewNumber++
			n.voteBlockHash = n.futureBlockHash
			n.futureBlockHash = ""
		} else {
			fmt.Println("debug: ", n.hashes[msg.VoteBlockHash])
		}

	case n.myViewNumber+1 == msg.ViewNumber: // this is proposal for new block with current sign
		// if not leader
		if !isLeader(msg.ViewNumber, msg.SenderID) {
			panic(fmt.Sprintf("expected to be leader, got %d", n.id))
		}

		if !msg.Block { // there is no block from future leader, it's incorrect behavior
			panic(fmt.Sprintf("expected block from future leader, got %v", msg))
		}
		if msg.Proof == nil {
			panic(fmt.Sprintf("expected proof, got %v", msg))
		}
		n.blockchain.AddBlock(msg.Proof.BlockHash)
		n.futureBlockHash = msg.VoteBlockHash
		n.myViewNumber++

		n.vote(msg.VoteBlockHash)

	default:
		panic(fmt.Sprintf("expected view number %d, got %d", n.myViewNumber, msg.ViewNumber))
	}
}

func proposeNewBlock(number ViewNumber) Block {
	return Block{number + 1}
}

func (n *Node) checkViewNumberBelongsToNodeSender(msg Msg) {
	if msg.ViewNumber%NumNodes != msg.SenderID {
		panic(fmt.Sprintf("expected view number %d to belong to node %d, got %d", msg.ViewNumber, n.id, msg.ViewNumber%NumNodes))
	}
}

func (n *Node) isLeader(number ViewNumber) bool {
	return isLeader(number, n.id)
}

func isLeader(number ViewNumber, id int) bool {
	return number%NumNodes == id

}

func (n *Node) vote(msg BlockHash) VotePower {
	n.hashes[msg]++
	return n.hashes[msg]
}

func (n *Node) SoftTimeout(msg Msg) {

}

func (n *Node) Timeout() {

}

func (n *Node) Tick() {

}

func channels(i int) []chan Msg {
	ch := make([]chan Msg, i)
	for i := range ch {
		ch[i] = make(chan Msg, 1)
	}
	return ch
}
