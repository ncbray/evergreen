package runtime

import (
	"strconv"
)

// TODO flow type?

const (
	NORMAL = iota
	FAIL
	EXCEPTION
)

type State struct {
	Stream         []rune
	Index          int
	Flow           int
	LookaheadLevel int
	deepest        int
	Offset         int
}

func (state *State) Checkpoint() int {
	return state.Index + state.Offset
}

func (state *State) Recover(index int) {
	state.Index = index - state.Offset
	state.Flow = NORMAL
}

func (state *State) Read() (r rune) {
	if state.Index < len(state.Stream) {
		r = state.Stream[state.Index]
		state.Index += 1
	} else {
		state.Fail()
	}
	return
}

func (state *State) Peek() (r rune) {
	if state.Index < len(state.Stream) {
		return state.Stream[state.Index]
	} else {
		state.Fail()
		return 0
	}
}

func (state *State) Consume() {
	state.Index += 1
}

func (state *State) Slice(start int) string {
	return string(state.Stream[start-state.Offset : state.Index])
}

func (state *State) Fail() {
	if state.Index > state.deepest && state.LookaheadLevel == 0 {
		state.deepest = state.Index
	}
	state.Flow = FAIL
}

func (state *State) LookaheadBegin() int {
	state.LookaheadLevel += 1
	return state.Index
}

func (state *State) LookaheadNormal(index int) {
	state.LookaheadLevel -= 1
	state.Index = index
	state.Flow = NORMAL
}

func (state *State) LookaheadFail(index int) {
	state.LookaheadLevel -= 1
	state.Index = index
	state.Fail()
}

func (state *State) Deepest() int {
	return state.deepest + state.Offset
}

func (state *State) RuneName(pos int) string {
	pos -= state.Offset
	var name string
	if pos < len(state.Stream) {
		name = strconv.QuoteRune(state.Stream[pos])
	} else {
		name = "EOF"
	}
	return name
}

// Utility function for generated tests
func MakeState(input string) *State {
	return &State{Stream: []rune(input)}
}
