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
	assertInt(len(result.Filters), 2, t)
	{
		f := result.Filters[0]
		assertRune('a', f.Min, t)
		assertRune('z', f.Max, t)
	}
	{
		f := result.Filters[1]
		assertRune('_', f.Min, t)
		assertRune('_', f.Max, t)
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
