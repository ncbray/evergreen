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

func parseFunction(state *dub.DubState) *FuncDecl {
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
	return &FuncDecl{Name: name, ReturnTypes: returnTypes, Block: block}
}

func parseImplements(state *dub.DubState) ASTTypeRef {
	checkpoint := state.Checkpoint()
	getKeyword(state, "implements")
	if state.Flow != 0 {
		state.Recover(checkpoint)
		return nil
	}
	return dubx.ParseTypeRef(state)
}

func parseStructure(state *dub.DubState) *StructDecl {
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

	fields := []*FieldDecl{}
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
		fields = append(fields, &FieldDecl{Name: name, Type: t})
	}

	getPunc(state, "}")
	if state.Flow != 0 {
		return nil
	}

	return &StructDecl{
		Name:       name,
		Implements: implements,
		Fields:     fields,
	}
}

func parseTest(state *dub.DubState) *Test {
	rule := dubx.Ident(state)
	if state.Flow != 0 {
		return nil
	}
	name := dubx.Ident(state)
	if state.Flow != 0 {
		return nil
	}
	input := dubx.DecodeString(state)
	if state.Flow != 0 {
		return nil
	}
	dubx.S(state)
	destructure := dubx.ParseDestructure(state)
	if state.Flow != 0 {
		return nil
	}
	return &Test{Name: name, Rule: rule, Input: input, Destructure: destructure}
}

func parseFile(state *dub.DubState) *File {
	decls := []Decl{}
	tests := []*Test{}
	for {
		checkpoint := state.Checkpoint()
		text := dubx.Ident(state)
		if state.Flow != 0 {
			state.Recover(checkpoint)
			goto done
		}
		switch text {
		case "func":
			f := parseFunction(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				goto done
			}
			decls = append(decls, f)
		case "struct":
			f := parseStructure(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				goto done
			}
			decls = append(decls, f)
		case "test":
			t := parseTest(state)
			if state.Flow != 0 {
				state.Recover(checkpoint)
				goto done
			}
			tests = append(tests, t)
		default:
			panic("foo")
			state.Recover(checkpoint)
			goto done
		}
	}
done:
	if state.Index != len(state.Stream) {
		state.Fail()
		return nil
	}
	return &File{
		Decls: decls,
		Tests: tests,
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
	text := string(stream[start:end])
	fmt.Printf("%s: Unexpected @ %d:%d\n%s", filename, line, col, text)
	// TODO tabs?
	arrow := []rune{}
	for i := 0; i < col; i++ {
		arrow = append(arrow, ' ')
	}
	arrow = append(arrow, '^')
	fmt.Println(string(arrow))
}

func ParseDASM(filename string) *File {
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
