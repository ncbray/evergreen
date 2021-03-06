package tree

import (
	"evergreen/compiler"
	"evergreen/dub/runtime"
	"fmt"
)

func ParseDub(data []byte, offset int, status compiler.TaskStatus) *File {
	status.Begin()
	defer status.End()

	stream := []rune(string(data))
	state := &runtime.State{Stream: stream, Offset: offset}
	f := ParseFile(state)
	if state.Flow == 0 {
		return f
	} else {
		pos := state.Deepest()
		name := state.RuneName(pos)
		status.LocationError(pos, fmt.Sprintf("Unexpected %s", name))
		return nil
	}
}
