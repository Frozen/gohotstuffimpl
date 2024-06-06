package chain

import "errors"

type Tally struct {
	m map[int]struct{}
}

func NewTally() Tally {
	return Tally{m: make(map[int]struct{})}
}

func (t *Tally) Add(i UniqueID) error {
	if _, ok := t.m[i]; ok {
		return errors.New("already tallied")
	}
	t.m[i] = struct{}{}
	return nil
}

func (t *Tally) Len() int {
	return len(t.m)
}
