package dub

type Flow int

const (
	NORMAL Flow = iota
	FAIL
)

type DubState struct {
  stream []rune
  index int
  flow Flow
}

func (state *DubState) Read() rune {
  if state.index < len(state.stream) {
		temp := state.stream[state.index]
		state.index += 1
		return temp
  } else {
		state.flow = FAIL
		return 0
	}
}

func (state *DubState) Save() int {
	return state.index;
}

func (state *DubState) Restore(pos int) {
  state.index = pos
  state.flow = NORMAL
}


func (state *DubState) Reject() {
  state.flow = FAIL
}
