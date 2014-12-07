// Package assert implements testing helpers.
package assert

import (
	"testing"
)

func IntEquals(t *testing.T, actual int, expected int) {
	if actual != expected {
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

func StringEquals(t *testing.T, actual string, expected string) {
	if actual != expected {
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

func IntListEquals(t *testing.T, actualList []int, expectedList []int) {
	IntEquals(t, len(actualList), len(expectedList))
	for i, expected := range expectedList {
		IntEquals(t, actualList[i], expected)
	}
}

func IntListListEquals(t *testing.T, actualList [][]int, expectedList [][]int) {
	if len(actualList) != len(expectedList) {
		t.Fatalf("%#v != %#v", actualList, expectedList)
	}
	for i, expected := range expectedList {
		IntListEquals(t, actualList[i], expected)
	}
}
