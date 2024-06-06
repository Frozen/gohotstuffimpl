package chain

import "errors"

type Tally struct {
	m map[Vote]struct{}
}

func NewTally() Tally {
	return Tally{m: make(map[Vote]struct{})}
}

type Vote struct {
	Type       MessageType
	ViewNumber ViewNumber
	Node       UniqueID
}

func (t *Tally) Add(i Vote) error {
	if _, ok := t.m[i]; ok {
		return errors.New("already tallied")
	}
	t.m[i] = struct{}{}
	return nil
}

func (t *Tally) Len() int {
	return len(t.m)
}
