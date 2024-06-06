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

func (t *Tally) Len(p MessageType, number ViewNumber) int {
	i := 0
	for v, _ := range t.m {
		if p == v.Type && number == v.ViewNumber {
			i++
		}
	}
	return i
}
