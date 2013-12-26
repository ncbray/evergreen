package gentests

import (
	"evergreen/dub"
	"generated/dubx"
	"testing"
)

func TestRuneMatchOK(t *testing.T) {
	state := &dub.DubState{Stream: []rune("[a-z_]")}
	result := dubx.MatchRune(state)
	assertState(state, 6, 0, t)
	if len(result.Filters) != 0 {
		t.Errorf("Expected %d filters, got %d", 0, len(result.Filters))
	}
}

func TestRuneMatchBad(t *testing.T) {
	state := &dub.DubState{Stream: []rune("[")}
	result := dubx.MatchRune(state)
	assertState(state, 1, 1, t)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}
