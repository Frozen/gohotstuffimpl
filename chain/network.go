package chain

type Network interface {
	Broadcast(msg NetworkMsg)
	ReceiveCh(id UniqueID) chan NetworkMsg
}

type MessageType string

var (
	NewView   MessageType
	Prepare   MessageType
	PreCommit MessageType
	Commit    MessageType
	Decide    MessageType
)

type NetworkMsg struct {
	Type       MessageType
	Node       int
	ViewNumber int
	Signatures int
	Payload    []byte
}

type Broadcaster struct {
	Chains []chan NetworkMsg
}

func NewBroadcaster(chs []chan NetworkMsg) *Broadcaster {
	return &Broadcaster{Chains: chs}
}

func (a *Broadcaster) Broadcast(m NetworkMsg) {
	for _, c := range a.Chains {
		c <- m
	}
}

func (a *Broadcaster) ReceiveCh(i UniqueID) chan NetworkMsg {
	return a.Chains[i]
}

func channels(i int) []chan NetworkMsg {
	ch := make([]chan NetworkMsg, i)
	for i := range ch {
		ch[i] = make(chan NetworkMsg, 1)
	}
	return ch

}
