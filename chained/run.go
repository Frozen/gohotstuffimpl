package chained

import "context"

func Run() {
	n0 := NewNode(0)
	n1 := NewNode(1)
	ctx := context.Background()
	//n3 := NewNode(2)

	chs := channels(NumNodes)
	responder := NewResponder(chs)

	msg := n0.Genesis()

	go NewNodeRunner(n0, chs).Run(ctx, responder)
	go NewNodeRunner(n1, chs).Run(ctx, responder)

	n1.Apply(msg, responder)
	//n3.Apply(msg)

	select {}
}
