package text

import (
	"evergreen/assert"
	"testing"
)

func TestSimple(t *testing.T) {
	b, w := BufferedCodeWriter()
	w.Line("foo")
	w.Line("bar")
	assert.StringEquals(t, b.String(), "foo\nbar\n")
}

func TestEmpty(t *testing.T) {
	b, w := BufferedCodeWriter()
	w.Line("foo")
	w.EmptyLines(2)
	w.EmptyLines(1)
	w.Line("bar")
	w.EmptyLines(2)
	assert.StringEquals(t, b.String(), "foo\n\n\nbar\n")
}

func TestMargin(t *testing.T) {
	b, w := BufferedCodeWriter()
	w.AppendMargin("  ")
	w.Line("foo")
	w.AppendMargin("  ")
	w.Line("bar")
	w.RestoreMargin()
	w.Line("baz")
	w.RestoreMargin()
	w.Line("fiz")
	assert.StringEquals(t, b.String(), "  foo\n    bar\n  baz\nfiz\n")
}

func TestTrimmedMargin1(t *testing.T) {
	b, w := BufferedCodeWriter()
	w.AppendMargin("  ")
	w.Line("foo ")
	w.EmptyLines(1)
	w.Line("bar  ")
	w.RestoreMargin()
	assert.StringEquals(t, b.String(), "  foo\n\n  bar\n")
}

func TestTrimmedMargin2(t *testing.T) {
	b, w := BufferedCodeWriter()
	w.AppendMargin("# ")
	w.Line("foo ")
	w.EmptyLines(1)
	w.Line("bar  ")
	w.RestoreMargin()
	assert.StringEquals(t, b.String(), "# foo\n#\n# bar\n")
}
