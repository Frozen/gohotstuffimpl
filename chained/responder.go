package chained

type ResponderImpl struct {
	channels []chan Msg
}

func NewResponder(chs []chan Msg) *ResponderImpl {
	return &ResponderImpl{channels: chs}
}

func (a *ResponderImpl) Broadcast(m Msg) {
	for i, c := range a.channels {
		if i != m.SenderID {
			c <- m
		}
	}
}

func (a *ResponderImpl) Send(id int, m Msg) {
	a.channels[id] <- m
}
