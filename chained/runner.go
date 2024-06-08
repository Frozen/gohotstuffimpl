package chained

import (
	"context"
	"fmt"
	"time"
)

type NodeRunner struct {
	node *Node
	chs  []chan Msg
}

func NewNodeRunner(node *Node, chs []chan Msg) *NodeRunner {
	return &NodeRunner{
		node: node,
		chs:  chs,
	}
}

func (nr *NodeRunner) Run(ctx context.Context, r Responder) {
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			nr.node.Tick()
		case m := <-nr.chs[nr.node.id]:
			fmt.Println("v: ", m.ViewNumber)
			nr.node.Apply(m, r)
		}
	}
}

func NodeRun(ctx context.Context, node ...*Node) {
	for _, n := range node {
		go func(node *Node) {

		}(n)
	}
}
