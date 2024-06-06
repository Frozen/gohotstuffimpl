package chain

import (
	"context"
)

const NumNodes = 4

type ViewNumber = int
type UniqueID = int

type Chain struct {
	viewNumber   int
	Uniq         UniqueID
	network      Network
	tally        Tally
	viewNumberCh chan ViewNumber
}

func NewChain(uniq int, network Network) Chain {
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
				//c.n.Broadcast(NetworkMsg{Type: NewView})
			} else {
				//fmt.Println("sending new view")
				c.network.Broadcast(
					NetworkMsg{
						Type:       NewView,
						Node:       c.Uniq,
						ViewNumber: i,
					},
				)
			}
		case m := <-c.network.ReceiveCh(c.Uniq):

			switch m.Type {
			case NewView:
				if c.isLeader() {
					// i'm the leader
					err := c.tally.Add(m.Node)
					if err != nil {
						panic(err)
					}
					if c.tally.Len() == NumNodes {
						c.network.Broadcast(NetworkMsg{
							Type:       Prepare,
							Node:       c.Uniq,
							ViewNumber: c.viewNumber,
							Signatures: c.tally.Len(),
							Payload:    createProposal(),
						})
					}

				} else {

					/* Pre commit phase */

					// i'm not the leader
					c.viewNumberCh <- m.ViewNumber
				}

			}
		}

	}

}

func (c *Chain) isLeader() bool {
	return c.viewNumber%NumNodes == c.Uniq
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
