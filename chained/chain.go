package chained

import (
	"fmt"
	"strconv"
)

const NumNodes = 3

type Node struct {
	id            int
	myViewNumber  ViewNumber
	hashes        map[BlockHash]VotePower
	voteBlockHash BlockHash
}

func NewNode(id int) *Node {
	return &Node{
		id:           id,
		myViewNumber: 0,
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

	if msg.Block { // message from future leader
		if n.isLeader(n.myViewNumber + 1) { // broadcast new block
			Msg{
				Type:          Generic,
				ViewNumber:    msg.ViewNumber + 1,
				SenderID:      n.id,
				VoteBlockHash: strconv.Itoa(msg.ViewNumber),
				Block:         true,
			}
			return
		} else { //sign
			Msg{
				Type:          Generic,
				ViewNumber:    msg.ViewNumber,
				SenderID:      n.id,
				VoteBlockHash: strconv.Itoa(msg.ViewNumber),
				Block:         false,
			}
			return
		}
	} else { // message from validator
		// i should be leader
		switch {
		case n.myViewNumber == msg.ViewNumber: // here we collect sigs
			if !n.isLeader(n.myViewNumber) {
				fmt.Println("just skip")
				return
			}
			if n.vote(msg.VoteBlockHash) == NumNodes { // 100%
				responder.Broadcast(Msg{
					Type:          Generic,
					ViewNumber:    msg.ViewNumber + 1,
					SenderID:      n.id,
					Block:         false,
					VoteBlockHash: strconv.Itoa(msg.ViewNumber),
					Proof: &Proof{
						BlockHash: msg.VoteBlockHash,
						Proof:     nil,
					},
				})
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
			n.vote(msg.VoteBlockHash)

		default:
			panic(fmt.Sprintf("expected view number %d, got %d", n.myViewNumber, msg.ViewNumber))
		}

		Msg{
			Type: Generic,
		}
		return

	}

	//if msg.ViewNumber%NumNodes == n.id {
	//
	//}
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

func (n *Node) SoftTimer() {

}

func main() {
	n1 := NewNode(0)
	n2 := NewNode(1)
	n3 := NewNode(2)

	msg := n1.Msg(0)
	n2.Apply(msg)
	n3.Apply(msg)

}
