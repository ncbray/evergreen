package base

import (
	"bytes"
	"testing"
)

func bufferedWriter() (*bytes.Buffer, *CodeWriter) {
	b := &bytes.Buffer{}
	w := &CodeWriter{Out: b}
	return b, w
}

func checkString(actual string, expected string, t *testing.T) {
	if actual != expected {
		t.Fatalf("%#v != %#v", actual, expected)
	}
}

func TestSimple(t *testing.T) {
	b, w := bufferedWriter()
	w.Line("foo")
	w.Line("bar")
	checkString(b.String(), "foo\nbar\n", t)
}

func TestEmpty(t *testing.T) {
	b, w := bufferedWriter()
	w.Line("foo")
	w.EmptyLines(2)
	w.EmptyLines(1)
	w.Line("bar")
	w.EmptyLines(2)
	checkString(b.String(), "foo\n\n\nbar\n", t)
}

func TestMargin(t *testing.T) {
	b, w := bufferedWriter()
	w.PushMargin("  ")
	w.Line("foo")
	w.PushMargin("  ")
	w.Line("bar")
	w.PopMargin()
	w.Line("baz")
	w.PopMargin()
	w.Line("fiz")
	checkString(b.String(), "  foo\n    bar\n  baz\nfiz\n", t)
}

func TestTrimmedMargin1(t *testing.T) {
	b, w := bufferedWriter()
	w.PushMargin("  ")
	w.Line("foo ")
	w.EmptyLines(1)
	w.Line("bar  ")
	w.PopMargin()
	checkString(b.String(), "  foo\n\n  bar\n", t)
}

func TestTrimmedMargin2(t *testing.T) {
	b, w := bufferedWriter()
	w.PushMargin("# ")
	w.Line("foo ")
	w.EmptyLines(1)
	w.Line("bar  ")
	w.PopMargin()
	checkString(b.String(), "# foo\n#\n# bar\n", t)
}
