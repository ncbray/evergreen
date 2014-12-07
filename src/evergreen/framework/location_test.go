package framework

import (
	"evergreen/assert"
	"testing"
)

func checkLocation(stream []rune, lines []int, pos int, eLine int, eCol int, eText string, t *testing.T) {
	line, col, text := getLocation(stream, lines, pos)
	assert.IntEquals(t, line, eLine)
	assert.IntEquals(t, col, eCol)
	assert.StringEquals(t, text, eText)
}

func TestFileLinesFull(t *testing.T) {
	stream := []rune("a\nb\nc")
	lines := findLines(stream)
	assert.IntListEquals(t, lines, []int{0, 2, 4})
	checkLocation(stream, lines, 0, 0, 0, "a", t)
	checkLocation(stream, lines, 1, 0, 1, "a", t)
	checkLocation(stream, lines, 2, 1, 0, "b", t)
	checkLocation(stream, lines, 3, 1, 1, "b", t)
	checkLocation(stream, lines, 4, 2, 0, "c", t)
	checkLocation(stream, lines, 5, 2, 1, "c", t)
	checkLocation(stream, lines, 6, 2, 2, "c", t)
}

func TestFileLinesEmpty(t *testing.T) {
	stream := []rune("\nx\n\n")
	lines := findLines(stream)
	assert.IntListEquals(t, lines, []int{0, 1, 3, 4})
	checkLocation(stream, lines, 0, 0, 0, "", t)
	checkLocation(stream, lines, 1, 1, 0, "x", t)
	checkLocation(stream, lines, 2, 1, 1, "x", t)
	checkLocation(stream, lines, 3, 2, 0, "", t)
	checkLocation(stream, lines, 4, 3, 0, "", t)
	checkLocation(stream, lines, 5, 3, 1, "", t)
}
