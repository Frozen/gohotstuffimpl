package chain

import (
	"context"
	"fmt"
)

func NewChain(uniq UniqueID, network Network) Chain {
	return Chain{
		Uniq:         uniq,
		network:      network,
		tally:        NewTally(),
		viewNumberCh: make(chan ViewNumber, 1),
	}
}

func (c *Chain) Run(ctx context.Context) {
	c.viewNumberCh <- 0
	for {
		select {
		case <-ctx.Done():
			return
		case i := <-c.viewNumberCh:
			c.viewNumber = i
			if c.isLeader() {
				// i'm the leader
				//c.n.Broadcast(Msg{Type: NewView})
			} else {
				//fmt.Println("sending new view")
				c.network.Broadcast(
					Msg{
						Type:       NewView,
						Node:       c.Uniq,
						ViewNumber: i,
					},
				)
			}
		case m := <-c.network.ReceiveCh(c.Uniq):
			if c.isLeader() {
				switch m.Type {
				case NewView: // leader: collect new view, send prepare
					err := c.tally.Add(Vote{Type: NewView, ViewNumber: m.ViewNumber, Node: m.Node})
					if err != nil {
						panic(err)
					}
					if c.tally.Len(NewView, m.ViewNumber) == NumNodes {
						c.network.Broadcast(Msg{
							Justify: &QC{
								Type:       Prepare,
								Node:       c.Uniq,
								ViewNumber: c.viewNumber,
							},
							Payload: createProposal(),
						})
					}
				case Prepare: // leader: collect prepare, send pre-commit

					err := c.tally.Add(Vote{
						Type:       Prepare,
						ViewNumber: m.ViewNumber,
						Node:       m.Node,
					})
					if err != nil {
						panic(err)
					}
					// if we collected enough prepare votes
					if c.tally.Len(Prepare, m.ViewNumber) == NumNodes {
						c.prepareQC = &QC{
							Type:       Prepare,
							Node:       c.Uniq,
							ViewNumber: c.viewNumber,
						}
						c.network.Broadcast(Msg{
							Type:       PreCommit,
							ViewNumber: m.ViewNumber,
							Node:       c.Uniq,
							Payload:    createProposal(),
							Justify:    c.prepareQC,
						})
					}
				case PreCommit: // leader: pre-commit -> commit
					err := c.tally.Add(Vote{
						Type:       PreCommit,
						ViewNumber: m.ViewNumber,
					}
				}
			} else { // validator
				switch m.Type {
				case Prepare: //
					if m.Justify == nil { // validator: prepare -> pre-commit
						if !isLeaderForView(m.ViewNumber, m.Node) {
							panic(fmt.Sprintf("invalid message received to node %d: %v", c.Uniq, m))
						}
						c.network.Broadcast(Msg{
							Type:       Prepare,
							ViewNumber: m.ViewNumber,
							Node:       c.Uniq,
							Justify:    nil,
							Payload:    nil,
						})
					} else { // validator: prepare -> pre-commit
						if !isLeaderForView(m.ViewNumber, m.Node) {
							panic("invalid message received")
						}
						c.prepareQC = m.Justify
						c.network.Broadcast(m)
					}

				case PreCommit: // validator: collect pre-commit, send commit
					c.lockedQC = m.Justify
				}
			}

		}

	}

}

func (c *Chain) isLeader() bool {
	return isLeaderForView(c.viewNumber, c.Uniq)
}

func isLeaderForView(viewNumber int, uniq UniqueID) bool {
	return UniqueID(viewNumber%NumNodes) == uniq
}

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	chains := Chains{}
	chains.Run(ctx)
	select {}
}

func createProposal() []byte {
	return []byte("proposal")
}
