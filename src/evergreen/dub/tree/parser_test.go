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
		offset := p.AddFile(filename, []rune(string(data)))
		status := framework.MakeStatus(p)
		ParseDub(data, offset, status)
	}
}

func BenchmarkSemantic(b *testing.B) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		b.Fatal(err)
	}
	p := framework.MakeProvider()
	offset := p.AddFile(filename, []rune(string(data)))
	status := framework.MakeStatus(p)
	file := ParseDub(data, offset, status)
	pkg := &Package{
		Files: []*File{
			file,
		},
	}
	index := MakeBuiltinTypeIndex()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		glbls := MakeDubGlobals(index)
		SemanticPass(pkg, glbls, status.CreateChild())
	}
}
