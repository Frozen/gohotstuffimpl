package chain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTally(t *testing.T) {
	tally := NewTally()
	if len(tally.m) != 0 {
		t.Errorf("expected 0, got %d", len(tally.m))
	}
}

func TestTallyAdd(t *testing.T) {
	tally := NewTally()
	err := tally.Add(Vote{Prepare, 0, 0})
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	err = tally.Add(Vote{Prepare, 0, 0})
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestTallyLen(t *testing.T) {
	tally := NewTally()
	require.Equal(t, 0, tally.Len(Prepare, 0))

	//if  != 0 {
	//	t.Errorf("expected 0, got %d", tally.Len())
	//}
	tally.Add(Vote{Prepare, 0, 0})
	require.Equal(t, 1, tally.Len(Prepare, 0))
	//tally.Add(1)
	//if tally.Len() != 1 {
	//	t.Errorf("expected 1, got %d", tally.Len())
	//}
}
