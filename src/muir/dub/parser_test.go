package dub

import (
	"testing"
)

func assertMatched(state *DubState, result *RuneMatch, invert bool, filters int, t *testing.T) {
	if state.flow != NORMAL {
		t.Errorf("Expected flow would be %d, but it was %d", NORMAL, state.flow)
	}
	if state.index != len(state.stream) {
		t.Errorf("Expected index would be %d, but it was %d", len(state.stream), state.index)
	}
	if result == nil {
		t.Errorf("Expected non-nil result.")
	}
	if result.Invert != invert {
		t.Errorf("Expected %t invert, got %t", invert, result.Invert)
	}
	if len(result.Filters) != filters {
		t.Errorf("Expected %d filters, got %d", filters, len(result.Filters))
	}
}

func assertFailed(state *DubState, pos int, result *RuneMatch, t *testing.T) {
	if state.flow != FAIL {
		t.Errorf("Expected flow would be %d, but it was %d", FAIL, state.flow)
	}
	if state.index != pos {
		t.Errorf("Expected index would be %d, but it was %d", pos, state.index)
	}
	if result != nil {
		t.Errorf("Expected nil result.")
	}
}

func TestBadStart(t *testing.T) {
	state := &DubState{stream: []rune("a-zA-Z_]")}
	result := Dub_RuneMatch(state)
	assertFailed(state, 1, result, t)
}

func TestBadEnd(t *testing.T) {
	state := &DubState{stream: []rune("[")}
	result := Dub_RuneMatch(state)
	assertFailed(state, 1, result, t)
}

func TestBadRange(t *testing.T) {
	state := &DubState{stream: []rune("[a-]")}
	result := Dub_RuneMatch(state)
	assertFailed(state, 4, result, t)
}

func TestSimple(t *testing.T) {
	state := &DubState{stream: []rune("[0]")}
	result := Dub_RuneMatch(state)
	assertMatched(state, result, false, 1, t)
}

func TestSimpleInvert(t *testing.T) {
	state := &DubState{stream: []rune("[^0]")}
	result := Dub_RuneMatch(state)
	assertMatched(state, result, true, 1, t)
}

func TestRangeInvert(t *testing.T) {
	state := &DubState{stream: []rune("[^a-z]")}
	result := Dub_RuneMatch(state)
	assertMatched(state, result, true, 1, t)
}

func TestIdentifier(t *testing.T) {
	state := &DubState{stream: []rune("[a-zA-Z_]")}
	result := Dub_RuneMatch(state)
	assertMatched(state, result, false, 3, t)
}

func TestEscape(t *testing.T) {
	state := &DubState{stream: []rune(`[ \t\r\n]`)}
	result := Dub_RuneMatch(state)
	assertMatched(state, result, false, 4, t)
}
