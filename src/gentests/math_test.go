package gentests

import (
	"evergreen/dub"
	"generated/math"
	"testing"
)

func assertState(state *dub.DubState, index int, flow int, t *testing.T) {
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

func TestDigits(t *testing.T) {
	state := &dub.DubState{Stream: []rune("123  4")}
	math.Digits(state)
	assertState(state, 3, 0, t)
	math.S(state)
	assertState(state, 5, 0, t)
	math.Digits(state)
	assertState(state, 6, 0, t)
	math.S(state)
	assertState(state, 6, 0, t)
	math.Digits(state)
	assertState(state, 6, 1, t)
}

func TestInteger(t *testing.T) {
	state := &dub.DubState{Stream: []rune("123  4")}
	result := math.Integer(state)
	assertState(state, 5, 0, t)
	assertString("123", result.Text, t)

	result = math.Integer(state)
	assertState(state, 6, 0, t)
	assertString("4", result.Text, t)

	result = math.Integer(state)
	assertState(state, 6, 1, t)
	if result != nil {
		t.Errorf("Expected %#v, got %#v", nil, result)
	}
}

func TestAddOK(t *testing.T) {
	state := &dub.DubState{Stream: []rune("1  + 234 ")}
	result := math.Add(state)
	assertState(state, 9, 0, t)
	assertString("1", result.Left.Text, t)
	assertString("+", result.Op, t)
	assertString("234", result.Right.Text, t)
}

func TestAddBad(t *testing.T) {
	state := &dub.DubState{Stream: []rune("1 234")}
	result := math.Add(state)
	assertState(state, 3, 1, t)
	if result != nil {
		t.Errorf("Expected %#v, got %#v", nil, result)
	}
}