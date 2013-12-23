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
