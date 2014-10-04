package tree

import (
	"evergreen/framework"
	"io/ioutil"
	"testing"
)

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
