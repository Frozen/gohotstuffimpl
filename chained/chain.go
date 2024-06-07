package chained

import (
	"fmt"
	"strconv"
)

func NewNode(id int) *Node {
	return &Node{
		id:           id,
		myViewNumber: 0,
		blockchain:   NewBlockchain(),
	}
}

func (n *Node) Msg(viewNumber ViewNumber) Msg {
	return Msg{
		Type:          Generic,
		ViewNumber:    viewNumber,
		VoteBlockHash: strconv.Itoa(viewNumber - 1),
	}
}

func (n *Node) Apply(msg Msg, responder Responder) {
	switch { // I should be leader
	case n.myViewNumber == msg.ViewNumber: // here we collect sigs
		if msg.VoteBlockHash != n.voteBlockHash {
			panic(fmt.Sprintf("expected vote block hash %s, got %s", n.voteBlockHash, msg.VoteBlockHash))
		}
		// TODO: check time as well

		if n.vote(msg.VoteBlockHash) == NumNodes { // 100%
			n.blockchain.AddBlock(msg.VoteBlockHash)

			responder.Broadcast(Msg{
				Type:          Generic,
				ViewNumber:    msg.ViewNumber + 1,
				SenderID:      n.id,
				Block:         false,
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

	case n.myViewNumber+1 == msg.ViewNumber: // this is next leader
		// if not leader
		if !isLeader(msg.ViewNumber, msg.SenderID) {
			panic(fmt.Sprintf("expected to be leader, got %d", n.id))
		}

		if !msg.Block { // there is no block from future leader, it's incorrect behavior
			panic(fmt.Sprintf("expected block from future leader, got %v", msg))
		}
		n.futureBlockHash = msg.VoteBlockHash
		n.vote(msg.VoteBlockHash)

	default:
		panic(fmt.Sprintf("expected view number %d, got %d", n.myViewNumber, msg.ViewNumber))
	}
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

func main() {
	n1 := NewNode(0)
	n2 := NewNode(1)
	n3 := NewNode(2)

	msg := n1.Msg(0)
	n2.Apply(msg)
	n3.Apply(msg)

}
