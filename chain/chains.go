package chain

import "context"

type Chains struct {
	Chains []Chain
}

func (c *Chains) Run(ctx context.Context) {
	chs := channels(NumNodes)
	broadcaster := NewBroadcaster(chs)
	for i := 0; i < NumNodes; i++ {
		c.Chains = append(c.Chains, NewChain(UniqueID(i), broadcaster))
	}
	for _, chain := range c.Chains {
		go chain.Run(ctx)
	}
}
