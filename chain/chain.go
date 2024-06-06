package chain

import (
	"context"
)

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
				case NewView:
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
				case Prepare:
					err := c.tally.Add(Vote{
						Type:       Prepare,
						ViewNumber: m.ViewNumber,
						Node:       m.Node,
					})
					if err != nil {
						panic(err)
					}
					if c.tally.Len(Prepare, m.ViewNumber) == NumNodes {
						c.network.Broadcast(Msg{
							Justify: prepareQc,
						})
					}
				}
			} else {
				switch m.Type {
				case Prepare:
					if m.QC == nil { // view change

					} else {

					}
					if !isLeaderForView(m.Node, m.ViewNumber) {
						panic("invalid message received")
					}
					c.network.Broadcast(m)
				}
			}

		}

	}

}

func (c *Chain) isLeader() bool {
	return c.isLeaderForView(c.viewNumber)
}

func isLeaderForView(viewNumber int, uniq UniqueID) bool {
	return viewNumber%NumNodes == uniq
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
