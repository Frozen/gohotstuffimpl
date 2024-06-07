package chain

type Network interface {
	Broadcast(msg Msg)
	ReceiveCh(id UniqueID) chan Msg
}

type MessageType string

var (
	NewView   MessageType
	Prepare   MessageType
	PreCommit MessageType
	Commit    MessageType
	Decide    MessageType
)

type Broadcaster struct {
	Chains []chan Msg
}

func NewBroadcaster(chs []chan Msg) *Broadcaster {
	return &Broadcaster{Chains: chs}
}

func (a *Broadcaster) Broadcast(m Msg) {
	for i, c := range a.Chains {
		if UniqueID(i) != m.Node {
			c <- m
		}
	}
}

func (a *Broadcaster) ReceiveCh(i UniqueID) chan Msg {
	return a.Chains[i]
}

func channels(i int) []chan Msg {
	ch := make([]chan Msg, i)
	for i := range ch {
		ch[i] = make(chan Msg, 1)
	}
	return ch

}
