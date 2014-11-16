package tree

import (
	"evergreen/framework"
	"testing"
)

var root = "../../../../dub"

func BenchmarkParser(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := framework.MakeProvider()
		status := framework.MakeStatus(p)
		ParseProgram(status, p, root)
	}
}

/*
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
		program := MakeProgramScope(index)
		SemanticPass(program, pkg, status.CreateChild())
	}
}
*/
