package test

import (
	"evergreen/assert"
	"evergreen/dub/runtime"
	"generated/playground"
	"testing"
)

func TestAdd(t *testing.T) {
	state := &runtime.State{}
	result := playground.Add(state, 3, 4)
	assert.IntEquals(t, state.Flow, runtime.NORMAL)
	assert.IntEquals(t, result, 7)
}
