package framework

import (
	"fmt"
	"testing"
)

func checkString(actual string, expected string, t *testing.T) {
	if actual != expected {
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

func checkInt(name string, actual int, expected int, t *testing.T) {
	if actual != expected {
		t.Fatalf("%s: %d != %d", name, actual, expected)
	}
}

func checkIntList(actualList []int, expectedList []int, t *testing.T) {
	checkInt("len", len(actualList), len(expectedList), t)
	for i, expected := range expectedList {
		checkInt(fmt.Sprint(i), actualList[i], expected, t)
	}
}

func checkLocation(stream []rune, lines []int, pos int, eLine int, eCol int, eText string, t *testing.T) {
	line, col, text := GetLocation(stream, lines, pos)
	checkInt("line", line, eLine, t)
	checkInt("col", col, eCol, t)
	checkString(text, eText, t)
}

func TestFileLinesFull(t *testing.T) {
	stream := []rune("a\nb\nc")
	lines := FindLines(stream)
	checkIntList(lines, []int{0, 2, 4}, t)
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
	lines := FindLines(stream)
	checkIntList(lines, []int{0, 1, 3, 4}, t)
	checkLocation(stream, lines, 0, 0, 0, "", t)
	checkLocation(stream, lines, 1, 1, 0, "x", t)
	checkLocation(stream, lines, 2, 1, 1, "x", t)
	checkLocation(stream, lines, 3, 2, 0, "", t)
	checkLocation(stream, lines, 4, 3, 0, "", t)
	checkLocation(stream, lines, 5, 3, 1, "", t)
}
