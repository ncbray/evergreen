package tree

import (
	"evergreen/framework"
	"fmt"
	"io/ioutil"
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

var filename = "../../../../dub/dub.dub"

func BenchmarkParser(b *testing.B) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := framework.MakeProvider()
		p.AddFile(filename, []rune(string(data)))
		status := framework.MakeStatus(p)
		ParseDub(data, status)
	}
}

func BenchmarkSemantic(b *testing.B) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		b.Fatal(err)
	}
	p := framework.MakeProvider()
	p.AddFile(filename, []rune(string(data)))
	status := framework.MakeStatus(p)
	file := ParseDub(data, status)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		SemanticPass(file, status.CreateChild())
	}
}
