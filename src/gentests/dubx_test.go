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

func TestSequence(t *testing.T) {
	state := &dub.DubState{Stream: []rune("[1-2] [3][4-5]")}
	result := dubx.Sequence(state)
	assertState(state, 14, 0, t)
	s, ok := result.(*dubx.MatchSequence)
	if !ok {
		t.Errorf("Not MatchSequence: %v", result)
	}
	assertInt(len(s.Matches), 3, t)
	{
		m, ok := s.Matches[0].(*dubx.RuneRangeMatch)
		if !ok {
			t.Errorf("Not RuneMatch: %v", m)
		}
		assertInt(len(m.Filters), 1, t)
		f := m.Filters[0]
		assertRune('1', f.Min, t)
		assertRune('2', f.Max, t)
	}
	{
		m, ok := s.Matches[1].(*dubx.RuneRangeMatch)
		if !ok {
			t.Errorf("Not RuneMatch: %v", m)
		}
		assertInt(len(m.Filters), 1, t)
		f := m.Filters[0]
		assertRune('3', f.Min, t)
		assertRune('3', f.Max, t)
	}
	{
		m, ok := s.Matches[2].(*dubx.RuneRangeMatch)
		if !ok {
			t.Errorf("Not RuneMatch: %v", m)
		}
		assertInt(len(m.Filters), 1, t)
		f := m.Filters[0]
		assertRune('4', f.Min, t)
		assertRune('5', f.Max, t)
	}
}
