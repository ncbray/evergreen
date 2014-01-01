package gentests

import (
	"evergreen/dub/runtime"
	"generated/dub/tree"
	"testing"
)

func assertState(state *runtime.State, index int, flow int, t *testing.T) {
	if state.Index != index {
		t.Errorf("Expected index %d, got %d", index, state.Index)
	}
	if state.Flow != flow {
		t.Errorf("Expected flow %d, got %d", flow, state.Flow)
	}
}

func assertString(expected string, actual string, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func assertInt(expected int, actual int, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected %#v, got %#v", expected, actual)
	}
}

func assertRune(expected rune, actual rune, t *testing.T) {
	if actual != expected {
		t.Errorf("Expected %#U, got %#U", expected, actual)
	}
}

func TestRuneMatchOK(t *testing.T) {
	state := &runtime.State{Stream: []rune("[a-z_]")}
	result := tree.MatchRune(state)
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
	state := &runtime.State{Stream: []rune("[")}
	result := tree.MatchRune(state)
	assertState(state, 1, 1, t)
	if result != nil {
		t.Errorf("Expected nil, got %v", result)
	}
}

func TestSequence(t *testing.T) {
	state := &runtime.State{Stream: []rune("[1-2] [3][4-5]")}
	result := tree.Sequence(state)
	assertState(state, 14, 0, t)
	s, ok := result.(*tree.MatchSequence)
	if !ok {
		t.Errorf("Not MatchSequence: %v", result)
	}
	assertInt(len(s.Matches), 3, t)
	{
		m, ok := s.Matches[0].(*tree.RuneRangeMatch)
		if !ok {
			t.Errorf("Not RuneMatch: %v", m)
		}
		assertInt(len(m.Filters), 1, t)
		f := m.Filters[0]
		assertRune('1', f.Min, t)
		assertRune('2', f.Max, t)
	}
	{
		m, ok := s.Matches[1].(*tree.RuneRangeMatch)
		if !ok {
			t.Errorf("Not RuneMatch: %v", m)
		}
		assertInt(len(m.Filters), 1, t)
		f := m.Filters[0]
		assertRune('3', f.Min, t)
		assertRune('3', f.Max, t)
	}
	{
		m, ok := s.Matches[2].(*tree.RuneRangeMatch)
		if !ok {
			t.Errorf("Not RuneMatch: %v", m)
		}
		assertInt(len(m.Filters), 1, t)
		f := m.Filters[0]
		assertRune('4', f.Min, t)
		assertRune('5', f.Max, t)
	}
}
