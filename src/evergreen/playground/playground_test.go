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

func TestCoerceFraction(t *testing.T) {
	state := &runtime.State{}
	result := playground.CoerceFraction(state)
	assert.IntEquals(t, state.Flow, runtime.NORMAL)
	assert.Float32Equals(t, result, 0.625)
}

func TestBlend(t *testing.T) {
	state := &runtime.State{}
	result := playground.Blend(state, 3, 11, 0.25)
	assert.IntEquals(t, state.Flow, runtime.NORMAL)
	assert.Float32Equals(t, result, 5)
}

func TestFooProxy(t *testing.T) {
	state := &runtime.State{}
	result := playground.FooProxy(state)
	assert.IntEquals(t, state.Flow, runtime.NORMAL)
	assert.IntEquals(t, result, 37)
}

func TestExplicitSpecialization(t *testing.T) {
	state := &runtime.State{}
	result := playground.ExplicitSpecialization(state)
	assert.IntEquals(t, state.Flow, runtime.NORMAL)
	assert.IntListEquals(t, result, []int{1})
}

func TestStringAddition(t *testing.T) {
	state := &runtime.State{}
	result := playground.StringAddition(state)
	assert.IntEquals(t, state.Flow, runtime.NORMAL)
	assert.StringEquals(t, result, "foobar")
}
