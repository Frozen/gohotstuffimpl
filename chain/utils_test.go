package chain

import "testing"

func TestNewTally(t *testing.T) {
	tally := NewTally()
	if len(tally.m) != 0 {
		t.Errorf("expected 0, got %d", len(tally.m))
	}
}

func TestTallyAdd(t *testing.T) {
	tally := NewTally()
	err := tally.Add(1)
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
	err = tally.Add(1)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestTallyLen(t *testing.T) {
	tally := NewTally()
	if tally.Len() != 0 {
		t.Errorf("expected 0, got %d", tally.Len())
	}
	tally.Add(1)
	if tally.Len() != 1 {
		t.Errorf("expected 1, got %d", tally.Len())
	}
}
