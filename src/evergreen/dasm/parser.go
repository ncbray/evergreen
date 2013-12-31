package dasm

import (
	"evergreen/dub"
	"evergreen/dubx"
	"fmt"
	"io/ioutil"
)

func getPunc(state *dub.DubState, value string) {
	for _, expected := range []rune(value) {
		c := state.Read()
		if state.Flow != 0 || c != expected {
			state.Fail()
			return
		}
	}
	dubx.S(state)
}

func getKeyword(state *dub.DubState, expected string) {
	for _, expected := range []rune(expected) {
		actual := state.Peek()
		if state.Flow != 0 {
			return
		}
		if actual != expected {
			state.Fail()
			return
		}
		state.Consume()
	}
	dubx.S(state)
}

func parseFunction(state *dub.DubState) *dubx.FuncDecl {
	getKeyword(state, "func")
	if state.Flow != 0 {
		return nil
	}
	name := dubx.Ident(state)
	if state.Flow != 0 {
		return nil
	}
	returnTypes := dubx.ParseTypeList(state)
	if state.Flow != 0 {
		return nil
	}
	block := dubx.ParseCodeBlock(state)
	if state.Flow != 0 {
		return nil
	}
	return &dubx.FuncDecl{Name: name, ReturnTypes: returnTypes, Block: block}
}

func parseImplements(state *dub.DubState) dubx.ASTTypeRef {
	checkpoint := state.Checkpoint()
	getKeyword(state, "implements")
	if state.Flow != 0 {
		state.Recover(checkpoint)
		return nil
	}
	return dubx.ParseTypeRef(state)
}

func parseStructure(state *dub.DubState) *dubx.StructDecl {
	getKeyword(state, "struct")
	if state.Flow != 0 {
		return nil
	}
	name := dubx.Ident(state)
	if state.Flow != 0 {
		return nil
	}

	implements := parseImplements(state)
	if state.Flow != 0 {
		return nil
	}

	getPunc(state, "{")
	if state.Flow != 0 {
		return nil
	}

	fields := []*dubx.FieldDecl{}
	for {
		checkpoint := state.Checkpoint()
		name := dubx.Ident(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		t := dubx.ParseTypeRef(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			break
		}
		fields = append(fields, &dubx.FieldDecl{Name: name, Type: t})
	}

	getPunc(state, "}")
	if state.Flow != 0 {
		return nil
	}

	return &dubx.StructDecl{
		Name:       name,
		Implements: implements,
		Fields:     fields,
	}
}

func parseFile(state *dub.DubState) *dubx.File {
	decls := []dubx.ASTDecl{}
	tests := []*dubx.Test{}
	for {
		checkpoint := state.Checkpoint()
		f := parseFunction(state)
		if state.Flow == 0 {
			decls = append(decls, f)
			continue
		}
		state.Recover(checkpoint)
		s := parseStructure(state)
		if state.Flow == 0 {
			decls = append(decls, s)
			continue
		}
		state.Recover(checkpoint)
		t := dubx.ParseTest(state)
		if state.Flow == 0 {
			tests = append(tests, t)
			continue
		}
		state.Recover(checkpoint)
		break
	}
	// Fail if not all input was consumed.
	if state.Flow == 0 && state.Index != len(state.Stream) {
		state.Fail()
		return nil
	}
	return &dubx.File{
		Decls: decls,
		Tests: tests,
	}
}

func ResolveType(ref dubx.ASTTypeRef) dubx.ASTType {
	switch ref := ref.(type) {
	case *dubx.TypeRef:
		return ref.T
	case *dubx.ListTypeRef:
		return ref.T
	default:
		panic(ref)
	}
}

func FindLines(stream []rune) []int {
	lines := []int{}
	for i, r := range stream {
		if r == '\n' {
			lines = append(lines, i+1)
		}
	}
	lines = append(lines, len(stream))
	return lines
}

func PrintError(filename string, deepest int, stream []rune, lines []int) {
	// Stupid linear search
	var line int
	var col int
	var start int
	var end int
	for i, s := range lines {
		start = end
		end = s
		line = i + 1
		if deepest < end {
			break
		}
	}
	col = deepest - start

	// HACK trim newline
	for end > start && (end > len(stream) || stream[end-1] == '\n' || stream[end-1] == '\t') {
		end -= 1
	}
	text := string(stream[start:end])
	fmt.Printf("%s: Unexpected @ %d:%d\n%s\n", filename, line, col, text)
	// TODO tabs?
	arrow := []rune{}
	for i := 0; i < col; i++ {
		arrow = append(arrow, ' ')
	}
	arrow = append(arrow, '^')
	fmt.Println(string(arrow))
}

func ParseDASM(filename string) *dubx.File {
	data, _ := ioutil.ReadFile(filename)
	stream := []rune(string(data))
	state := &dub.DubState{Stream: stream}
	f := parseFile(state)
	if state.Flow != 0 {
		lines := FindLines(stream)
		PrintError(filename, state.Deepest, stream, lines)
		panic(state.Deepest)
	}
	return f
}
